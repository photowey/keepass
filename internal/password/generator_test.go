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

package password

import (
	"strings"
	"testing"
)

func TestAlphabetForPreset(t *testing.T) {
	alphabet, err := AlphabetForPreset(PresetSymbols)
	if err != nil {
		t.Fatalf("AlphabetForPreset() error = %v", err)
	}

	if !strings.ContainsRune(alphabet, '!') || !strings.ContainsRune(alphabet, '+') {
		t.Fatalf("expected symbols preset to include special symbols, got %q", alphabet)
	}

	if _, err := AlphabetForPreset("unknown"); err == nil {
		t.Fatal("expected error for unsupported preset")
	}
}

func TestGenerateUsesRequestedAlphabetAndLength(t *testing.T) {
	alphabet := "abC123!@#-_"
	value, err := Generate(32, alphabet)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if len(value) != 32 {
		t.Fatalf("expected length 32, got %d", len(value))
	}

	for _, r := range value {
		if !strings.ContainsRune(alphabet, r) {
			t.Fatalf("generated rune %q not found in alphabet", r)
		}
	}
}

func TestGenerateRejectsInvalidInput(t *testing.T) {
	if _, err := Generate(0, "abc"); err == nil {
		t.Fatal("expected error for zero length")
	}

	if _, err := Generate(8, ""); err == nil {
		t.Fatal("expected error for empty alphabet")
	}
}
