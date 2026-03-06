package holons

import "testing"

func TestNormalizeBuildTargetSupportsAppleFamilies(t *testing.T) {
	tests := map[string]string{
		"macos":    "macos",
		"ios":      "ios",
		"tvos":     "tvos",
		"watchos":  "watchos",
		"visionos": "visionos",
	}

	for input, want := range tests {
		got, err := normalizeBuildTarget(input)
		if err != nil {
			t.Fatalf("normalizeBuildTarget(%q) returned error: %v", input, err)
		}
		if got != want {
			t.Fatalf("normalizeBuildTarget(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestIsValidPlatformAcceptsAppleFamilies(t *testing.T) {
	for _, platform := range []string{"macos", "ios", "tvos", "watchos", "visionos"} {
		if !isValidPlatform(platform) {
			t.Fatalf("expected %q to be a valid platform", platform)
		}
	}
}
