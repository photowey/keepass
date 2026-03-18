package credentialaudit_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/photowey/keepass/internal/credentialaudit"
	"github.com/photowey/keepass/internal/vault"
)

func ExampleAnalyze() {
	now := time.Date(2026, time.March, 18, 12, 0, 0, 0, time.UTC)
	report := credentialaudit.Analyze([]vault.Entry{
		{
			Alias:             "github",
			Username:          "alice",
			Password:          "shared-secret",
			URI:               "https://github.com",
			PasswordUpdatedAt: now.AddDate(0, 0, -10),
		},
		{
			Alias:             "gitlab",
			Password:          "shared-secret",
			PasswordUpdatedAt: now,
		},
	}, 7, now)

	fmt.Println(report.OverallStatus)
	for _, finding := range report.Findings {
		fmt.Printf("%s %s\n", finding.Type, strings.Join(finding.Aliases, ","))
	}

	// Output:
	// warn
	// stale_password github
	// missing_username gitlab
	// missing_uri gitlab
	// duplicate_password github,gitlab
}
