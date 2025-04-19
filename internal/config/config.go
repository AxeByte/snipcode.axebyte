package config

import (
  "errors"
  "os"

  "github.com/urfave/cli/v2"
  "gopkg.in/yaml.v3"
)

type Config struct {
  DefaultName    string   `yaml:"default_name"`
  IgnorePatterns []string `yaml:"ignore_patterns"`
}

// InitConfig writes a default .grepattern.yaml
func InitConfig(c *cli.Context) error {
  cfg := Config{DefaultName: "comp_code.txt", IgnorePatterns: []string{}}
  data, err := yaml.Marshal(cfg)
  if err != nil {
    return err
  }
  return os.WriteFile(".grepattern.yaml", data, 0644)
}

// LoadConfig loads and parses .grepattern.yaml
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
