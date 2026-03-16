package api

import "strings"

const defaultVersion = "v0.5.0"

var (
	Version = defaultVersion
	Commit  = "unknown"
)

func VersionString() string {
	version := strings.TrimSpace(Version)
	if version == "" {
		version = defaultVersion
	}

	commit := strings.TrimSpace(Commit)
	if commit == "" || commit == "unknown" {
		return version
	}
	if len(commit) > 7 {
		commit = commit[:7]
	}
	return version + " (" + commit + ")"
}

func Banner() string {
	return "op " + VersionString()
}
