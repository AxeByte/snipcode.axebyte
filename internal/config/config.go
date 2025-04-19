package config

import (
	"errors"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DefaultName    string   `yaml:"default_name"`
	IgnorePatterns []string `yaml:"ignore_patterns"`
}

// common defaults for media and VCS dirs
var defaultIgnores = []string{
	".git/**",
	"*.png", "*.jpg", "*.jpeg", "*.gif", "*.webp", "*.svg", "*.ico",
	"*.mp4", "*.mov", "*.avi", "*.mkv",
}

// InitConfig writes a default .grepattern.yaml to CWD.
func InitConfig(c *cli.Context) error {
	cfg := Config{
		DefaultName:    "comp_code.txt",
		IgnorePatterns: append([]string{}, defaultIgnores...),
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := os.WriteFile(".grepattern.yaml", data, 0644); err != nil {
		return err
	}
	log.Info("Created local config .grepattern.yaml")
	return nil
}

// InitAdmin writes a default .grepattern.yaml under $XDG_CONFIG_HOME/snipcode/, with exhaustive patterns.
func InitAdmin(c *cli.Context) error {
	patterns := []string{
		"dist/**", ".cache/**", "examples/**", "helm-chart/*", ".yarn/*", "yarn.lock",
		".cursor/**", "seeds.go", "*/dist/*", "*/node_modules/*", "tests/*.go",
		".next/**", "static/**", ".DS_Store", "*.toml", "LICENSE", "docs/**",
		"*.sqlite", "*.lock", "*.ink", "*.lockb", "*.test.*", "*.css", "*.jpeg",
		"docs.go", "logs/**", "deploy.sh", "lefthook.*", ".gitignore", ".env",
		"*_test*", "*dock*", "*Dock*", "images/*", "*.g4", "*txt*", "output.txt",
		"README.md", "aaa.json", ".github/**", "package-lock.json", "migrations/**",
		"venv/**", "__pycache__/**", "go.mod", "go.sum", "Dockerfile",
	}
	// merge with media + .git defaults
	patterns = append(patterns, defaultIgnores...)

	cfg := Config{
		DefaultName:    "comp_code.txt",
		IgnorePatterns: patterns,
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	targetDir := filepath.Join(configDir, "snipcode")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(targetDir, ".grepattern.yaml")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	log.Infof("Created global config at %s", path)
	return nil
}

// LoadConfig loads and parses a .grepattern.yaml from the given path.
func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("could not read config: " + err.Error())
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, errors.New("invalid config: " + err.Error())
	}
	return &cfg, nil
}
