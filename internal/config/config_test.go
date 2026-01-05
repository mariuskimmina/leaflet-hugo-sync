package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	content := `
source:
  handle: "test.bsky.social"
  collection: "com.whtwnd.blog.entry"
output:
  posts_dir: "content/posts"
  images_dir: "static/images"
  image_path_prefix: "/images"
template:
  frontmatter: |
    ---
    title: "{{ .Title }}"
    ---
`
	tmpfile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Source.Handle != "test.bsky.social" {
		t.Errorf("expected test.bsky.social, got %s", cfg.Source.Handle)
	}
}
