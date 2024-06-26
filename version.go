package main

import (
	"fmt"

	goversion "github.com/hashicorp/go-version"
)

// To set this from outside, use go build -ldflags "-X 'main.version=$(VERSION)'"
var version string = "0.1.0"

// A pre-release marker for the version. If this is "" (empty string)
// then it means that it is a final release. Otherwise, this is a pre-release
// such as "dev" (in development), "beta", "rc1", etc.
var prerelease = ""

func init() {
	// Verify that the version is proper semantic version, which should always be the case.
	_, err := goversion.NewVersion(version)
	if err != nil {
		panic(err.Error())
	}
}

// VersionString returns the complete version string, including prerelease
func VersionString() string {
	if prerelease != "" {
		return fmt.Sprintf("%s-%s", version, prerelease)
	}
	return version
}
