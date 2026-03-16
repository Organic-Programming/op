package api_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func withWorkingDir(t *testing.T, dir string) {
	t.Helper()

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir(%s): %v", dir, err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(original)
	})
}

func monorepoRoot(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
}

func copyPlatformManifestProto(t *testing.T, root string) {
	t.Helper()

	source := filepath.Join(monorepoRoot(t), "_protos", "holons", "v1", "manifest.proto")
	data, err := os.ReadFile(source)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", source, err)
	}

	targetDir := filepath.Join(root, "_protos", "holons", "v1")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", targetDir, err)
	}
	if err := os.WriteFile(filepath.Join(targetDir, "manifest.proto"), data, 0o644); err != nil {
		t.Fatalf("WriteFile(manifest.proto): %v", err)
	}
}

func writeProtoHolon(t *testing.T, root string) string {
	t.Helper()

	copyPlatformManifestProto(t, root)

	dir := filepath.Join(root, "alpha-service")
	if err := os.MkdirAll(filepath.Join(dir, "api", "v1"), 0o755); err != nil {
		t.Fatalf("MkdirAll(api/v1): %v", err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "cmd", "alpha-service"), 0o755); err != nil {
		t.Fatalf("MkdirAll(cmd): %v", err)
	}

	goMod := "module example.com/alpha-service\n\ngo 1.24.0\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatalf("WriteFile(go.mod): %v", err)
	}

	mainGo := `package main

func main() {}
`
	if err := os.WriteFile(filepath.Join(dir, "cmd", "alpha-service", "main.go"), []byte(mainGo), 0o644); err != nil {
		t.Fatalf("WriteFile(main.go): %v", err)
	}

	holonProto := fmt.Sprintf(`syntax = "proto3";

package alpha.v1;

import "holons/v1/manifest.proto";

option go_package = "example.com/alpha-service/gen/go/alpha/v1;alphav1";

option (holons.v1.manifest) = {
  identity: {
    schema: "holon/v1"
    uuid: "11111111-1111-1111-1111-111111111111"
    given_name: "Alpha"
    family_name: "Service"
    motto: "Answers ping."
    composer: "Test Suite"
    status: "draft"
    born: "2026-03-16"
  }
  description: "Proto-first test holon."
  lang: "go"
  kind: "native"
  build: {
    runner: "go-module"
    main: "./cmd/alpha-service"
  }
  requires: {
    commands: ["go"]
    files: ["go.mod"]
  }
  artifacts: {
    binary: "alpha-service"
  }
  contract: {
    proto: "api/v1/holon.proto"
    service: "alpha.v1.AlphaService"
    rpcs: ["Ping"]
  }
  sequences: [{
    name: "greet"
    description: "Render a greeting step."
    params: [{
      name: "name"
      description: "Name to greet."
      default: "world"
    }]
    steps: ["echo hello {{.name}}"]
  }]
};

service AlphaService {
  rpc Ping (PingRequest) returns (PingResponse);
}

message PingRequest {
  string name = 1;
}

message PingResponse {
  string message = 1;
}
`)
	if err := os.WriteFile(filepath.Join(dir, "api", "v1", "holon.proto"), []byte(holonProto), 0o644); err != nil {
		t.Fatalf("WriteFile(holon.proto): %v", err)
	}

	return dir
}
