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

package vault

import (
	"os"
	"testing"
)

func TestInspectFileReadsV1Metadata(t *testing.T) {
	store, _ := newTestStore(t)

	meta, err := InspectFile(store.Path())
	if err != nil {
		t.Fatalf("InspectFile() error = %v", err)
	}

	if meta.FormatVersion != 1 {
		t.Fatalf("expected format version 1, got %d", meta.FormatVersion)
	}

	if meta.KDF != "argon2id" || meta.Cipher != "xchacha20poly1305" {
		t.Fatalf("unexpected metadata: %+v", meta)
	}

	if meta.Argon2id == nil || meta.Argon2id.MemoryKiB == 0 {
		t.Fatalf("expected argon2id metadata, got %+v", meta)
	}
}

func TestInspectFileRejectsMissingVault(t *testing.T) {
	if _, err := InspectFile("/tmp/does-not-exist-keepass.kp"); err == nil {
		t.Fatal("expected missing vault error")
	}
}

func TestInspectFileRejectsInvalidHeader(t *testing.T) {
	store, _ := newTestStore(t)

	if err := os.WriteFile(store.Path(), []byte("bad"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := InspectFile(store.Path()); err == nil {
		t.Fatal("expected inspect error for invalid vault")
	}
}
