package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Source   Source   `yaml:"source"`
	Output   Output   `yaml:"output"`
	Template Template `yaml:"template"`
}

type Source struct {
	Handle          string `yaml:"handle"`
	Collection      string `yaml:"collection"`
	PublicationName string `yaml:"publication_name"`
}

type Output struct {
	PostsDir        string `yaml:"posts_dir"`
	ImagesDir       string `yaml:"images_dir"`
	ImagePathPrefix string `yaml:"image_path_prefix"`
}

type Template struct {
	Frontmatter string `yaml:"frontmatter"`
	Content     string `yaml:"content"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
