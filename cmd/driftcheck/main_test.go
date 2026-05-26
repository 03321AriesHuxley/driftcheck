package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestMainHelp verifies the binary can be built and prints usage without error.
func TestMainHelp(t *testing.T) {
	if os.Getenv("DRIFTCHECK_INTEGRATION") == "" {
		t.Skip("skipping integration test; set DRIFTCHECK_INTEGRATION=1 to run")
	}

	bin := buildBinary(t)
	cmd := exec.Command(bin, "--help")
	out, err := cmd.CombinedOutput()
	// --help exits with code 2 via flag package
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if !ok || exitErr.ExitCode() != 2 {
			t.Fatalf("unexpected error running --help: %v\noutput: %s", err, out)
		}
	}
	outStr := string(out)
	for _, expected := range []string{"-file", "-service", "-json"} {
		if !containsStr(outStr, expected) {
			t.Errorf("expected %q in help output, got:\n%s", expected, outStr)
		}
	}
}

// TestMainMissingComposeFile checks that a missing compose file exits with code 1.
func TestMainMissingComposeFile(t *testing.T) {
	if os.Getenv("DRIFTCHECK_INTEGRATION") == "" {
		t.Skip("skipping integration test; set DRIFTCHECK_INTEGRATION=1 to run")
	}

	bin := buildBinary(t)
	cmd := exec.Command(bin, "-file", "/nonexistent/docker-compose.yml")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing compose file")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok || exitErr.ExitCode() != 1 {
		t.Fatalf("expected exit code 1, got: %v", err)
	}
}

func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "driftcheck")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = filepath.Join("..", "..", "cmd", "driftcheck")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstring(s, sub))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
