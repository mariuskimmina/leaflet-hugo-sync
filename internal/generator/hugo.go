package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"mariuskimmina.com/leaflet-hugo-sync/internal/config"
)

type Generator struct {
	Cfg *config.Config
}

type PostData struct {
	Title       string
	CreatedAt   string
	Slug        string
	Filename    string
	Handle      string
	OriginalURL string
	Content     string
	Data        map[string]interface{}
}

func NewGenerator(cfg *config.Config) *Generator {
	return &Generator{Cfg: cfg}
}

func (g *Generator) GeneratePost(data PostData) error {
	// 1. Generate Frontmatter
	tmplFM, err := template.New("frontmatter").Parse(g.Cfg.Template.Frontmatter)
	if err != nil {
		return err
	}

	var bufFM bytes.Buffer
	if err := tmplFM.Execute(&bufFM, data); err != nil {
		return err
	}

	// 2. Generate Content
	contentTmplStr := g.Cfg.Template.Content
	if contentTmplStr == "" {
		contentTmplStr = "{{ .Content }}" // Default
	}

	tmplContent, err := template.New("content").Parse(contentTmplStr)
	if err != nil {
		return err
	}

	var bufContent bytes.Buffer
	if err := tmplContent.Execute(&bufContent, data); err != nil {
		return err
	}

	if err := os.MkdirAll(g.Cfg.Output.PostsDir, 0755); err != nil {
		return err
	}

	// Use Filename if provided, otherwise fall back to Slug
	filename := data.Filename
	if filename == "" {
		filename = data.Slug
	}
	fileName := filename + ".md"
	filePath := filepath.Join(g.Cfg.Output.PostsDir, fileName)

	fullContent := bufFM.String() + "\n" + bufContent.String()

	return os.WriteFile(filePath, []byte(fullContent), 0644)
}
