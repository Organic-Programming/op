package identity

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveFromProtoFileParsesCopyArtifactStep(t *testing.T) {
	root := t.TempDir()
	sharedProto := filepath.Join("..", "..", "..", "..", "_protos", "holons", "v1", "manifest.proto")
	data, err := os.ReadFile(sharedProto)
	if err != nil {
		t.Fatalf("ReadFile(%q) failed: %v", sharedProto, err)
	}

	sharedDir := filepath.Join(root, "_protos", "holons", "v1")
	if err := os.MkdirAll(sharedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sharedDir, "manifest.proto"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	protoDir := filepath.Join(root, "app", "api", "v1")
	if err := os.MkdirAll(protoDir, 0o755); err != nil {
		t.Fatal(err)
	}
	protoPath := filepath.Join(protoDir, "holon.proto")
	proto := `syntax = "proto3";

package test.v1;

import "holons/v1/manifest.proto";

option (holons.v1.manifest) = {
  identity: {
    schema: "holon/v1"
    uuid: "11111111-2222-3333-4444-555555555555"
    given_name: "Proto"
    family_name: "Artifact"
    motto: "Parses copy_artifact steps."
    composer: "test"
    status: "draft"
    born: "2026-03-16"
  }
  kind: "composite"
  build: {
    runner: "recipe"
    members: { id: "daemon" path: "../daemon" type: "holon" }
    members: { id: "app" path: "." type: "component" }
    targets: {
      key: "macos"
      value: {
        steps: { build_member: "daemon" }
        steps: {
          copy_artifact: {
            from: "daemon"
            to: "build/MyApp.app/Contents/Resources/Holons/daemon.holon"
          }
        }
      }
    }
  }
  artifacts: {
    primary: "build/MyApp.app"
  }
};
`
	if err := os.WriteFile(protoPath, []byte(proto), 0o644); err != nil {
		t.Fatal(err)
	}

	resolved, err := ResolveFromProtoFile(protoPath)
	if err != nil {
		t.Fatalf("ResolveFromProtoFile failed: %v", err)
	}

	target := resolved.BuildTargets["macos"]
	if len(target.Steps) != 2 {
		t.Fatalf("len(target.Steps) = %d, want 2", len(target.Steps))
	}
	if target.Steps[1].CopyArtifact == nil {
		t.Fatal("expected copy_artifact step to be resolved")
	}
	if got := target.Steps[1].CopyArtifact.From; got != "daemon" {
		t.Fatalf("CopyArtifact.From = %q, want daemon", got)
	}
	if got := target.Steps[1].CopyArtifact.To; got != "build/MyApp.app/Contents/Resources/Holons/daemon.holon" {
		t.Fatalf("CopyArtifact.To = %q", got)
	}
}
