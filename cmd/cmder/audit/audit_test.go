package audit

import (
	"bytes"
	"strings"
	"testing"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/manager"
	"github.com/photowey/keepass/internal/testutil"
)

func TestNewRejectsUnexpectedArguments(t *testing.T) {
	cmd := New()
	cmd.SetArgs([]string{"unexpected"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected argument validation error")
	}
}

func TestNewPrintsCredentialAuditJSON(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	mgr := manager.New(env, cfg)
	if err := mgr.Initialize("master", false); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	if _, _, err := mgr.Add("master", manager.AddInput{
		Alias:    "github",
		Username: "octocat",
		Password: "secret-123",
		URI:      "https://github.com",
	}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	cmd := New()
	cmd.SetArgs([]string{"--json", "--max-password-age-days", "0"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, `"overall_status": "ok"`) || !strings.Contains(output, `"max_age_days": 0`) {
		t.Fatalf("unexpected audit output %q", output)
	}
}
