package version

import "testing"

func TestNowPrefixesReleaseVersion(t *testing.T) {
	originalVersion := version
	t.Cleanup(func() {
		version = originalVersion
	})

	version = "1.2.3"

	if got := Now(); got != "v1.2.3" {
		t.Fatalf("expected prefixed release version, got %q", got)
	}
}

func TestNowKeepsDevVersion(t *testing.T) {
	originalVersion := version
	t.Cleanup(func() {
		version = originalVersion
	})

	version = "dev"

	if got := Now(); got != "dev" {
		t.Fatalf("expected dev version, got %q", got)
	}
}

func TestSummaryIncludesMetadata(t *testing.T) {
	originalVersion := version
	originalCommit := commit
	originalBuildTime := buildTime
	t.Cleanup(func() {
		version = originalVersion
		commit = originalCommit
		buildTime = originalBuildTime
	})

	version = "v2.0.0"
	commit = "abc1234"
	buildTime = "2026-03-18T12:00:00Z"

	got := Summary()
	want := "v2.0.0 (commit abc1234, built 2026-03-18T12:00:00Z)"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
