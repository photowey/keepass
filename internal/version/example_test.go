package version

import "fmt"

func ExampleSummary() {
	originalVersion := version
	originalCommit := commit
	originalBuildTime := buildTime
	defer func() {
		version = originalVersion
		commit = originalCommit
		buildTime = originalBuildTime
	}()

	version = "1.2.3"
	commit = "abc1234"
	buildTime = "2026-03-18T12:00:00Z"

	fmt.Println(Now())
	fmt.Println(Summary())

	// Output:
	// v1.2.3
	// v1.2.3 (commit abc1234, built 2026-03-18T12:00:00Z)
}
