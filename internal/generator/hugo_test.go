package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marius/leaflet-hugo-sync/internal/config"
)

func TestGeneratePost(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "output")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		Output: config.Output{
			PostsDir: filepath.Join(tmpDir, "posts"),
		},
		Template: config.Template{
			Frontmatter: "---\ntitle: \"{{ .Title }}\"\n---",
		},
	}

	gen := NewGenerator(cfg)
	data := PostData{
		Title:   "Hello World",
		Slug:    "hello-world",
		Content: "This is a test post.",
	}

	if err := gen.GeneratePost(data); err != nil {
		t.Fatalf("GeneratePost failed: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "posts", "hello-world.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", expectedPath)
	}

	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatal(err)
	}

	expectedContent := "---\ntitle: \"Hello World\"\n---\nThis is a test post."
	if string(content) != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, string(content))
	}
}
