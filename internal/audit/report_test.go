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

package audit

import (
	"testing"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/testutil"
	"github.com/photowey/keepass/internal/vault"
)

func TestCollectReportsHealthyVault(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	store := vault.NewStore(cfg.ResolveVaultPath(env), cfg)
	if err := store.Initialize("master-password", false); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	report, err := Collect()
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	if report.Config.Present != true || report.Vault.Present != true {
		t.Fatalf("expected config and vault present, got %+v", report)
	}

	if report.RehashRecommended {
		t.Fatalf("expected no rehash recommendation, got %+v", report)
	}
}

func TestCollectRecommendsRehashWhenConfigDiffers(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	store := vault.NewStore(cfg.ResolveVaultPath(env), cfg)
	if err := store.Initialize("master-password", false); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	cfg.Security.Argon2id.MemoryKiB = 16 * 1024
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save(updated cfg) error = %v", err)
	}

	report, err := Collect()
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	if !report.RehashRecommended {
		t.Fatalf("expected rehash recommendation, got %+v", report)
	}
}
