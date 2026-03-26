/*
 * Copyright © 2023-present the keepass authors. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
