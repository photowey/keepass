package credentialaudit

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/photowey/keepass/internal/vault"
)

type FindingType string

const (
	FindingStalePassword     FindingType = "stale_password"
	FindingDuplicatePassword FindingType = "duplicate_password"
	FindingMissingUsername   FindingType = "missing_username"
	FindingMissingURI        FindingType = "missing_uri"
)

type Finding struct {
	Type    FindingType `json:"type"`
	Aliases []string    `json:"aliases"`
	Message string      `json:"message"`
}

type Report struct {
	OverallStatus string    `json:"overall_status"`
	MaxAgeDays    int       `json:"max_age_days"`
	Findings      []Finding `json:"findings"`
}

func Analyze(entries []vault.Entry, maxAgeDays int, now time.Time) Report {
	report := Report{
		OverallStatus: "ok",
		MaxAgeDays:    maxAgeDays,
	}

	if maxAgeDays > 0 {
		cutoff := now.UTC().AddDate(0, 0, -maxAgeDays)
		for _, entry := range entries {
			if entry.PasswordUpdatedAt.Before(cutoff) {
				report.addFinding(Finding{
					Type:    FindingStalePassword,
					Aliases: []string{entry.Alias},
					Message: "Password age exceeds the configured threshold",
				})
			}
		}
	}

	passwordGroups := map[string][]string{}
	for _, entry := range entries {
		if strings.TrimSpace(entry.Password) != "" {
			passwordGroups[entry.Password] = append(passwordGroups[entry.Password], entry.Alias)
		}

		if strings.TrimSpace(entry.Username) == "" {
			report.addFinding(Finding{
				Type:    FindingMissingUsername,
				Aliases: []string{entry.Alias},
				Message: "Entry is missing a username",
			})
		}

		if strings.TrimSpace(entry.URI) == "" {
			report.addFinding(Finding{
				Type:    FindingMissingURI,
				Aliases: []string{entry.Alias},
				Message: "Entry is missing a URI",
			})
		}
	}

	for _, aliases := range passwordGroups {
		if len(aliases) > 1 {
			sort.Strings(aliases)
			report.addFinding(Finding{
				Type:    FindingDuplicatePassword,
				Aliases: aliases,
				Message: "Entries share the same password",
			})
		}
	}

	return report
}

func (r *Report) addFinding(f Finding) {
	r.Findings = append(r.Findings, f)
	r.OverallStatus = "warn"
}

func (r Report) MarshalJSON() ([]byte, error) {
	type alias Report
	return json.Marshal(alias(r))
}
