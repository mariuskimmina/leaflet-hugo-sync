package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/marius/leaflet-hugo-sync/internal/config"
)

type Generator struct {
	Cfg *config.Config
}

type PostData struct {
	Title     string
	CreatedAt string
	Slug      string
	Handle    string
	Content   string
	Data      map[string]interface{}
}

func NewGenerator(cfg *config.Config) *Generator {
	return &Generator{Cfg: cfg}
}

func (g *Generator) GeneratePost(data PostData) error {
	tmpl, err := template.New("frontmatter").Parse(g.Cfg.Template.Frontmatter)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	if err := os.MkdirAll(g.Cfg.Output.PostsDir, 0755); err != nil {
		return err
	}

	fileName := data.Slug + ".md"
	filePath := filepath.Join(g.Cfg.Output.PostsDir, fileName)

	fullContent := buf.String() + "\n" + data.Content

	return os.WriteFile(filePath, []byte(fullContent), 0644)
}

