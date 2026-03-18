package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/pkg/files"
)

var ErrVaultNotInitialized = errors.New("vault not initialized, run `keepass init` first")

type Store struct {
	path string
	cfg  configs.Config
	now  func() time.Time
}

func NewStore(path string, cfg configs.Config) *Store {
	return &Store{
		path: path,
		cfg:  cfg,
		now:  time.Now,
	}
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) Exists() bool {
	return files.Exists(s.path)
}

func (s *Store) Initialize(masterPassword string, force bool) error {
	if masterPassword == "" {
		return errors.New("master password cannot be blank")
	}

	if s.Exists() && !force {
		return errors.New("vault already exists")
	}

	document := NewDocument(s.now())
	return s.Save(masterPassword, document)
}

func (s *Store) Load(masterPassword string) (*Document, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrVaultNotInitialized
		}

		return nil, fmt.Errorf("read vault: %w", err)
	}

	plaintext, err := Decode(data, masterPassword)
	if err != nil {
		return nil, err
	}

	var document Document
	if err := json.Unmarshal(plaintext, &document); err != nil {
		return nil, fmt.Errorf("decode vault payload: %w", err)
	}

	return &document, nil
}

func (s *Store) Save(masterPassword string, document *Document) error {
	if document == nil {
		return errors.New("vault document cannot be nil")
	}

	document.UpdatedAt = s.now().UTC()

	plaintext, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("encode vault payload: %w", err)
	}

	data, err := Encode(s.cfg.Vault.FormatVersion, plaintext, masterPassword, s.cfg)
	if err != nil {
		return err
	}

	if err := files.EnsureDir(filepath.Dir(s.path), 0o700); err != nil {
		return fmt.Errorf("ensure vault dir: %w", err)
	}

	if err := files.WriteFileAtomic(s.path, data, 0o600); err != nil {
		return fmt.Errorf("write vault: %w", err)
	}

	return nil
}
