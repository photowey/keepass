package credentialaudit

import (
	"testing"
	"time"

	"github.com/photowey/keepass/internal/vault"
)

func TestAnalyzeFindsStaleDuplicateAndMissingMetadata(t *testing.T) {
	now := time.Unix(1_700_000_000, 0).UTC()
	entries := []vault.Entry{
		{
			Alias:             "github",
			Username:          "hellopass",
			Password:          "shared-secret",
			PasswordUpdatedAt: now.AddDate(0, 0, -365),
		},
		{
			Alias:             "gitea",
			Username:          "hellopass",
			Password:          "shared-secret",
			URI:               "https://gitea.example.com",
			PasswordUpdatedAt: now,
		},
	}

	report := Analyze(entries, 180, now)
	if report.OverallStatus != "warn" {
		t.Fatalf("expected warn status, got %+v", report)
	}

	if len(report.Findings) < 3 {
		t.Fatalf("expected multiple findings, got %+v", report)
	}
}
