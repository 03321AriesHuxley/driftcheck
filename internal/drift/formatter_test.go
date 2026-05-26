package drift

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatResults_AllClean(t *testing.T) {
	results := []Result{
		{ServiceName: "web", ContainerName: "app_web_1", ContainerID: "abc"},
		{ServiceName: "db", ContainerName: "app_db_1", ContainerID: "def"},
	}

	var buf bytes.Buffer
	count, err := FormatResults(&buf, results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 drift services, got %d", count)
	}
	out := buf.String()
	if !strings.Contains(out, "[OK]") {
		t.Error("expected [OK] in output")
	}
}

func TestFormatResults_WithDrift(t *testing.T) {
	r := Result{ServiceName: "web", ContainerName: "app_web_1"}
	r.AddDrift(DriftKindImage, "image", "nginx:1.25", "nginx:1.21")
	r.AddDrift(DriftKindEnv, "DEBUG", "false", "true")

	var buf bytes.Buffer
	count, err := FormatResults(&buf, []Result{r})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 drift service, got %d", count)
	}
	out := buf.String()
	if !strings.Contains(out, "[DRIFT]") {
		t.Error("expected [DRIFT] tag in output")
	}
	if !strings.Contains(out, "[IMAGE]") {
		t.Error("expected [IMAGE] kind in output")
	}
	if !strings.Contains(out, "nginx:1.25") {
		t.Error("expected expected image in output")
	}
}

func TestFormatResults_EmptyActual(t *testing.T) {
	r := Result{ServiceName: "cache", ContainerName: "app_cache_1"}
	r.AddDrift(DriftKindPort, "6379/tcp", "6379:6379", "")

	var buf bytes.Buffer
	_, err := FormatResults(&buf, []Result{r})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "<none>") {
		t.Error("expected <none> placeholder for empty actual")
	}
}

func TestFormatResults_Empty(t *testing.T) {
	var buf bytes.Buffer
	count, err := FormatResults(&buf, []Result{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}
