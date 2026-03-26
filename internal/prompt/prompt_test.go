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

package prompt

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestAskTrimsWhitespace(t *testing.T) {
	var out bytes.Buffer
	p := New(strings.NewReader("  alice  \n"), &out)

	value, err := p.Ask("Username")
	if err != nil {
		t.Fatalf("Ask() error = %v", err)
	}

	if value != "alice" {
		t.Fatalf("expected trimmed value, got %q", value)
	}

	if got := out.String(); got != "Username: " {
		t.Fatalf("unexpected prompt output %q", got)
	}
}

func TestAskDefaultUsesFallbackOnBlankInput(t *testing.T) {
	var out bytes.Buffer
	p := New(strings.NewReader("\n"), &out)

	value, err := p.AskDefault("URI", "https://example.com")
	if err != nil {
		t.Fatalf("AskDefault() error = %v", err)
	}

	if value != "https://example.com" {
		t.Fatalf("expected default value, got %q", value)
	}

	if got := out.String(); got != "URI [https://example.com]: " {
		t.Fatalf("unexpected prompt output %q", got)
	}
}

func TestAskDefaultUsesEnteredValue(t *testing.T) {
	var out bytes.Buffer
	p := New(strings.NewReader("https://override.example.com\n"), &out)

	value, err := p.AskDefault("URI", "https://example.com")
	if err != nil {
		t.Fatalf("AskDefault() error = %v", err)
	}

	if value != "https://override.example.com" {
		t.Fatalf("expected entered value, got %q", value)
	}
}

func TestAskOptionalDelegatesToAsk(t *testing.T) {
	var out bytes.Buffer
	p := New(strings.NewReader("  notes  \n"), &out)

	value, err := p.AskOptional("Note")
	if err != nil {
		t.Fatalf("AskOptional() error = %v", err)
	}

	if value != "notes" {
		t.Fatalf("expected trimmed optional value, got %q", value)
	}
}

func TestAskAcceptsEOFWithoutTrailingNewline(t *testing.T) {
	var out bytes.Buffer
	p := New(strings.NewReader("alice"), &out)

	value, err := p.Ask("Username")
	if err != nil {
		t.Fatalf("Ask() error = %v", err)
	}

	if value != "alice" {
		t.Fatalf("expected EOF value, got %q", value)
	}
}

func TestAskSecretFallsBackToReaderWhenNotTerminal(t *testing.T) {
	var out bytes.Buffer
	p := New(strings.NewReader("  secret-value  \n"), &out)

	value, err := p.AskSecret("Master password")
	if err != nil {
		t.Fatalf("AskSecret() error = %v", err)
	}

	if value != "secret-value" {
		t.Fatalf("expected trimmed secret, got %q", value)
	}

	if got := out.String(); got != "Master password: " {
		t.Fatalf("unexpected prompt output %q", got)
	}
}

func TestAskSecretAcceptsEOFWithoutTrailingNewline(t *testing.T) {
	var out bytes.Buffer
	p := New(strings.NewReader("secret-value"), &out)

	value, err := p.AskSecret("Master password")
	if err != nil {
		t.Fatalf("AskSecret() error = %v", err)
	}

	if value != "secret-value" {
		t.Fatalf("expected EOF secret, got %q", value)
	}
}

func TestAskSecretWithConfirmationRejectsMismatch(t *testing.T) {
	var out bytes.Buffer
	p := New(strings.NewReader("first\nsecond\n"), &out)

	_, err := p.AskSecretWithConfirmation("Master password", "Confirm master password")
	if err == nil {
		t.Fatal("expected confirmation mismatch error")
	}

	if !strings.Contains(err.Error(), "does not match confirmation") {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestAskSecretWithConfirmationAcceptsMatch(t *testing.T) {
	var out bytes.Buffer
	p := New(strings.NewReader("same-secret\nsame-secret\n"), &out)

	value, err := p.AskSecretWithConfirmation("Master password", "Confirm master password")
	if err != nil {
		t.Fatalf("AskSecretWithConfirmation() error = %v", err)
	}

	if value != "same-secret" {
		t.Fatalf("expected confirmed secret, got %q", value)
	}
}

func TestConfirmHonorsDefaultAndYesInput(t *testing.T) {
	t.Run("blank uses default", func(t *testing.T) {
		var out bytes.Buffer
		p := New(strings.NewReader("\n"), &out)

		ok, err := p.Confirm("Delete entry github", true)
		if err != nil {
			t.Fatalf("Confirm() error = %v", err)
		}

		if !ok {
			t.Fatal("expected default yes result")
		}
	})

	t.Run("explicit yes accepted", func(t *testing.T) {
		var out bytes.Buffer
		p := New(strings.NewReader("y\n"), &out)

		ok, err := p.Confirm("Delete entry github", false)
		if err != nil {
			t.Fatalf("Confirm() error = %v", err)
		}

		if !ok {
			t.Fatal("expected yes result")
		}
	})

	t.Run("explicit no rejected", func(t *testing.T) {
		var out bytes.Buffer
		p := New(strings.NewReader("no\n"), &out)

		ok, err := p.Confirm("Delete entry github", true)
		if err != nil {
			t.Fatalf("Confirm() error = %v", err)
		}

		if ok {
			t.Fatal("expected no result")
		}
	})

	t.Run("blank uses default no", func(t *testing.T) {
		var out bytes.Buffer
		p := New(strings.NewReader("\n"), &out)

		ok, err := p.Confirm("Delete entry github", false)
		if err != nil {
			t.Fatalf("Confirm() error = %v", err)
		}

		if ok {
			t.Fatal("expected default no result")
		}
	})

	t.Run("explicit yes word accepted", func(t *testing.T) {
		var out bytes.Buffer
		p := New(strings.NewReader("yes\n"), &out)

		ok, err := p.Confirm("Delete entry github", false)
		if err != nil {
			t.Fatalf("Confirm() error = %v", err)
		}

		if !ok {
			t.Fatal("expected yes word result")
		}
	})
}

type errWriter struct{}

func (errWriter) Write(p []byte) (int, error) {
	return 0, io.ErrClosedPipe
}

type errReader struct{}

func (errReader) Read(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestAskPropagatesWriterError(t *testing.T) {
	p := New(strings.NewReader("ignored\n"), errWriter{})

	if _, err := p.Ask("Username"); err == nil {
		t.Fatal("expected writer error")
	}
}

func TestAskPropagatesReaderError(t *testing.T) {
	var out bytes.Buffer
	p := New(errReader{}, &out)

	if _, err := p.Ask("Username"); err == nil {
		t.Fatal("expected reader error")
	}
}

func TestAskDefaultPropagatesReaderError(t *testing.T) {
	var out bytes.Buffer
	p := New(errReader{}, &out)

	if _, err := p.AskDefault("Username", "alice"); err == nil {
		t.Fatal("expected reader error")
	}
}

func TestAskSecretPropagatesReaderError(t *testing.T) {
	var out bytes.Buffer
	p := New(errReader{}, &out)

	if _, err := p.AskSecret("Master password"); err == nil {
		t.Fatal("expected reader error")
	}
}
