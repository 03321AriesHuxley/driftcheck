package compose_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftcheck/internal/compose"
)

const sampleCompose = `
version: "3.9"
services:
  web:
    image: nginx:1.25
    ports:
      - "80:80"
    environment:
      ENV: production
    restart: always
  db:
    image: postgres:15
    environment:
      POSTGRES_DB: mydb
    volumes:
      - db-data:/var/lib/postgresql/data
`

func writeTempCompose(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "docker-compose.yml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp compose file: %v", err)
	}
	return path
}

func TestParseFile_ValidCompose(t *testing.T) {
	path := writeTempCompose(t, sampleCompose)
	cf, err := compose.ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cf.Services) != 2 {
		t.Errorf("expected 2 services, got %d", len(cf.Services))
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := compose.ParseFile("/nonexistent/path/docker-compose.yml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestGetService_Existing(t *testing.T) {
	path := writeTempCompose(t, sampleCompose)
	cf, _ := compose.ParseFile(path)

	svc, err := cf.GetService("web")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.Image != "nginx:1.25" {
		t.Errorf("expected image nginx:1.25, got %s", svc.Image)
	}
	if svc.Restart != "always" {
		t.Errorf("expected restart always, got %s", svc.Restart)
	}
}

func TestGetService_Missing(t *testing.T) {
	path := writeTempCompose(t, sampleCompose)
	cf, _ := compose.ParseFile(path)

	_, err := cf.GetService("cache")
	if err == nil {
		t.Fatal("expected error for missing service, got nil")
	}
}
