package manager

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/credentialaudit"
	"github.com/photowey/keepass/internal/home"
	"github.com/photowey/keepass/internal/password"
	"github.com/photowey/keepass/internal/transfer"
	"github.com/photowey/keepass/internal/vault"
)

var aliasPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)

type Manager struct {
	env   home.Environment
	cfg   configs.Config
	store *vault.Store
	now   func() time.Time
}

type AddInput struct {
	Alias            string
	Username         string
	Password         string
	URI              string
	Note             string
	Tags             []string
	GeneratePassword bool
}

type UpdateInput struct {
	Username         *string
	Password         *string
	URI              *string
	Note             *string
	Tags             *[]string
	GeneratePassword bool
}

type ListFilter struct {
	Query string
	Tags  []string
}

func LoadCurrent() (*Manager, error) {
	env, err := home.Detect()
	if err != nil {
		return nil, err
	}

	cfg, err := configs.Load(env)
	if err != nil {
		return nil, err
	}

	return New(env, cfg), nil
}

func LoadOrCreateCurrent() (*Manager, error) {
	env, err := home.Detect()
	if err != nil {
		return nil, err
	}

	cfg, err := configs.LoadOrCreate(env)
	if err != nil {
		return nil, err
	}

	return New(env, cfg), nil
}

func New(env home.Environment, cfg configs.Config) *Manager {
	return &Manager{
		env:   env,
		cfg:   cfg,
		store: vault.NewStore(cfg.ResolveVaultPath(env), cfg),
		now:   time.Now,
	}
}

func (m *Manager) Env() home.Environment {
	return m.env
}

func (m *Manager) Config() configs.Config {
	return m.cfg
}

func (m *Manager) VaultPath() string {
	return m.store.Path()
}

func (m *Manager) VaultExists() bool {
	return m.store.Exists()
}

func (m *Manager) Initialize(masterPassword string, force bool) error {
	return m.store.Initialize(masterPassword, force)
}

func (m *Manager) Rehash(masterPassword string) (int, error) {
	document, err := m.store.Load(masterPassword)
	if err != nil {
		return 0, err
	}

	if err := m.store.Save(masterPassword, document); err != nil {
		return 0, err
	}

	return len(document.Entries), nil
}

func (m *Manager) Export(masterPassword string) (transfer.ExportDocument, error) {
	document, err := m.store.Load(masterPassword)
	if err != nil {
		return transfer.ExportDocument{}, err
	}

	entries := append([]vault.Entry(nil), document.Entries...)
	sortEntries(entries)
	return transfer.ExportDocument{
		Version:    transfer.ExportVersion,
		ExportedAt: m.now().UTC(),
		Entries:    entries,
	}, nil
}

func (m *Manager) Import(masterPassword string, doc transfer.ExportDocument, strategy string) (transfer.ImportResult, error) {
	document, err := m.store.Load(masterPassword)
	if err != nil {
		return transfer.ImportResult{}, err
	}

	strategy, err = transfer.NormalizeConflictStrategy(strategy)
	if err != nil {
		return transfer.ImportResult{}, err
	}

	result := transfer.ImportResult{}
	for _, imported := range doc.Entries {
		normalizedAlias, err := normalizeAlias(imported.Alias)
		if err != nil {
			return transfer.ImportResult{}, err
		}

		username := strings.TrimSpace(imported.Username)
		if username == "" {
			return transfer.ImportResult{}, errors.New("username cannot be blank")
		}

		if strings.TrimSpace(imported.Password) == "" {
			return transfer.ImportResult{}, errors.New("password cannot be blank")
		}

		tags, err := normalizeTags(imported.Tags)
		if err != nil {
			return transfer.ImportResult{}, err
		}

		imported.Alias = normalizedAlias
		imported.Username = username
		imported.URI = strings.TrimSpace(imported.URI)
		imported.Note = strings.TrimSpace(imported.Note)
		imported.Tags = tags

		index, _, found := findExact(document.Entries, normalizedAlias)
		if !found {
			document.Entries = append(document.Entries, imported)
			result.Added++
			continue
		}

		switch strategy {
		case transfer.ConflictFail:
			return transfer.ImportResult{}, fmt.Errorf("alias %q already exists", normalizedAlias)
		case transfer.ConflictSkip:
			result.Skipped++
		case transfer.ConflictOverwrite:
			document.Entries[index] = imported
			result.Overwrote++
		}
	}

	sortEntries(document.Entries)
	if err := m.store.Save(masterPassword, document); err != nil {
		return transfer.ImportResult{}, err
	}

	return result, nil
}

func (m *Manager) CreateBackup(path string, force bool) (string, error) {
	return transfer.CreateBackupBundle(path, m.env, m.cfg, force, m.now())
}

func RestoreCurrent(path string, force bool) error {
	env, err := home.Detect()
	if err != nil {
		return err
	}

	return transfer.RestoreBackupBundle(path, env, force)
}

func (m *Manager) AuditCredentials(masterPassword string, maxAgeDays int) (credentialaudit.Report, error) {
	document, err := m.store.Load(masterPassword)
	if err != nil {
		return credentialaudit.Report{}, err
	}

	entries := append([]vault.Entry(nil), document.Entries...)
	sortEntries(entries)
	return credentialaudit.Analyze(entries, maxAgeDays, m.now().UTC()), nil
}

func (m *Manager) Rotate(masterPassword, alias string, passwordValue *string, generate bool) (vault.Entry, bool, error) {
	input := UpdateInput{
		Password:         passwordValue,
		GeneratePassword: generate,
	}

	return m.Update(masterPassword, alias, input)
}

func (m *Manager) Add(masterPassword string, input AddInput) (vault.Entry, bool, error) {
	document, err := m.store.Load(masterPassword)
	if err != nil {
		return vault.Entry{}, false, err
	}

	entry, generated, err := m.newEntry(input)
	if err != nil {
		return vault.Entry{}, false, err
	}

	if _, _, found := findExact(document.Entries, entry.Alias); found {
		return vault.Entry{}, false, fmt.Errorf("alias %q already exists", entry.Alias)
	}

	document.Entries = append(document.Entries, entry)
	sortEntries(document.Entries)

	if err := m.store.Save(masterPassword, document); err != nil {
		return vault.Entry{}, false, err
	}

	return entry, generated, nil
}

func (m *Manager) List(masterPassword string, filter ListFilter) ([]vault.Entry, error) {
	document, err := m.store.Load(masterPassword)
	if err != nil {
		return nil, err
	}

	filterTags, err := normalizeTags(filter.Tags)
	if err != nil {
		return nil, err
	}

	query := strings.ToLower(strings.TrimSpace(filter.Query))
	results := make([]vault.Entry, 0, len(document.Entries))
	for _, entry := range document.Entries {
		if !matchesTags(entry.Tags, filterTags) {
			continue
		}

		if query != "" && !matchesQuery(entry, query) {
			continue
		}

		results = append(results, entry)
	}

	sortEntries(results)
	return results, nil
}

func (m *Manager) Get(masterPassword, alias string) (vault.Entry, error) {
	document, err := m.store.Load(masterPassword)
	if err != nil {
		return vault.Entry{}, err
	}

	return resolveEntry(document.Entries, alias)
}

func (m *Manager) Update(masterPassword, alias string, input UpdateInput) (vault.Entry, bool, error) {
	document, err := m.store.Load(masterPassword)
	if err != nil {
		return vault.Entry{}, false, err
	}

	index, entry, err := resolveEntryWithIndex(document.Entries, alias)
	if err != nil {
		return vault.Entry{}, false, err
	}

	generated := false
	now := m.now().UTC()

	if input.Username != nil {
		username := strings.TrimSpace(*input.Username)
		if username == "" {
			return vault.Entry{}, false, errors.New("username cannot be blank")
		}
		entry.Username = username
	}

	if input.URI != nil {
		entry.URI = strings.TrimSpace(*input.URI)
	}

	if input.Note != nil {
		entry.Note = strings.TrimSpace(*input.Note)
	}

	if input.Tags != nil {
		tags, err := normalizeTags(*input.Tags)
		if err != nil {
			return vault.Entry{}, false, err
		}
		entry.Tags = tags
	}

	if input.GeneratePassword {
		alphabet, err := m.cfg.PasswordGenerator.EffectiveAlphabet()
		if err != nil {
			return vault.Entry{}, false, err
		}

		entry.Password, err = password.Generate(m.cfg.PasswordGenerator.DefaultLength, alphabet)
		if err != nil {
			return vault.Entry{}, false, err
		}
		entry.PasswordUpdatedAt = now
		generated = true
	} else if input.Password != nil {
		if strings.TrimSpace(*input.Password) == "" {
			return vault.Entry{}, false, errors.New("password cannot be blank")
		}
		entry.Password = *input.Password
		entry.PasswordUpdatedAt = now
	}

	entry.UpdatedAt = now
	document.Entries[index] = entry

	if err := m.store.Save(masterPassword, document); err != nil {
		return vault.Entry{}, false, err
	}

	return entry, generated, nil
}

func (m *Manager) Delete(masterPassword, alias string) (vault.Entry, error) {
	document, err := m.store.Load(masterPassword)
	if err != nil {
		return vault.Entry{}, err
	}

	index, entry, err := resolveEntryWithIndex(document.Entries, alias)
	if err != nil {
		return vault.Entry{}, err
	}

	document.Entries = append(document.Entries[:index], document.Entries[index+1:]...)

	if err := m.store.Save(masterPassword, document); err != nil {
		return vault.Entry{}, err
	}

	return entry, nil
}

func (m *Manager) newEntry(input AddInput) (vault.Entry, bool, error) {
	alias, err := normalizeAlias(input.Alias)
	if err != nil {
		return vault.Entry{}, false, err
	}

	username := strings.TrimSpace(input.Username)
	if username == "" {
		return vault.Entry{}, false, errors.New("username cannot be blank")
	}

	tags, err := normalizeTags(input.Tags)
	if err != nil {
		return vault.Entry{}, false, err
	}

	accountPassword := input.Password
	generated := false
	if input.GeneratePassword || strings.TrimSpace(accountPassword) == "" {
		alphabet, err := m.cfg.PasswordGenerator.EffectiveAlphabet()
		if err != nil {
			return vault.Entry{}, false, err
		}

		accountPassword, err = password.Generate(m.cfg.PasswordGenerator.DefaultLength, alphabet)
		if err != nil {
			return vault.Entry{}, false, err
		}
		generated = true
	}

	now := m.now().UTC()

	return vault.Entry{
		Alias:             alias,
		Username:          username,
		Password:          accountPassword,
		URI:               strings.TrimSpace(input.URI),
		Note:              strings.TrimSpace(input.Note),
		Tags:              tags,
		CreatedAt:         now,
		UpdatedAt:         now,
		PasswordUpdatedAt: now,
	}, generated, nil
}

func normalizeAlias(alias string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(alias))
	if normalized == "" {
		return "", errors.New("alias cannot be blank")
	}

	if !aliasPattern.MatchString(normalized) {
		return "", errors.New("alias must match [a-z0-9][a-z0-9._-]*")
	}

	return normalized, nil
}

func normalizeTags(tags []string) ([]string, error) {
	seen := map[string]struct{}{}
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		trimmed := strings.ToLower(strings.TrimSpace(tag))
		if trimmed == "" {
			continue
		}

		if !aliasPattern.MatchString(trimmed) {
			return nil, fmt.Errorf("tag %q must match [a-z0-9][a-z0-9._-]*", trimmed)
		}

		if _, ok := seen[trimmed]; ok {
			continue
		}

		seen[trimmed] = struct{}{}
		normalized = append(normalized, trimmed)
	}

	sort.Strings(normalized)
	return normalized, nil
}

func resolveEntry(entries []vault.Entry, alias string) (vault.Entry, error) {
	_, entry, err := resolveEntryWithIndex(entries, alias)
	return entry, err
}

func resolveEntryWithIndex(entries []vault.Entry, alias string) (int, vault.Entry, error) {
	normalized, err := normalizeAlias(alias)
	if err != nil {
		return -1, vault.Entry{}, err
	}

	if index, entry, found := findExact(entries, normalized); found {
		return index, entry, nil
	}

	matches := make([]int, 0, 2)
	for index, entry := range entries {
		if strings.HasPrefix(entry.Alias, normalized) {
			matches = append(matches, index)
		}
	}

	switch len(matches) {
	case 0:
		return -1, vault.Entry{}, fmt.Errorf("no entry matches alias %q", normalized)
	case 1:
		index := matches[0]
		return index, entries[index], nil
	default:
		aliases := make([]string, 0, len(matches))
		for _, index := range matches {
			aliases = append(aliases, entries[index].Alias)
		}
		sort.Strings(aliases)
		return -1, vault.Entry{}, fmt.Errorf("alias %q is ambiguous: %s", normalized, strings.Join(aliases, ", "))
	}
}

func findExact(entries []vault.Entry, alias string) (int, vault.Entry, bool) {
	for index, entry := range entries {
		if entry.Alias == alias {
			return index, entry, true
		}
	}

	return -1, vault.Entry{}, false
}

func matchesTags(entryTags, required []string) bool {
	if len(required) == 0 {
		return true
	}

	set := map[string]struct{}{}
	for _, tag := range entryTags {
		set[tag] = struct{}{}
	}

	for _, tag := range required {
		if _, ok := set[tag]; !ok {
			return false
		}
	}

	return true
}

func matchesQuery(entry vault.Entry, query string) bool {
	if strings.HasPrefix(entry.Alias, query) {
		return true
	}

	candidates := []string{entry.Username, entry.URI, entry.Note}
	candidates = append(candidates, entry.Tags...)
	for _, candidate := range candidates {
		if strings.Contains(strings.ToLower(candidate), query) {
			return true
		}
	}

	return false
}

func sortEntries(entries []vault.Entry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Alias < entries[j].Alias
	})
}
