package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version   string                   `yaml:"version"`
	Analyzers []*ProjectAnalyzerConfig `yaml:"analyzers"`
}

type ProjectAnalyzerConfig struct {
	Prompt   string    `yaml:"prompt"`
	Analyzer string    `yaml:"analyzer"`
	Model    string    `yaml:"model"`
	Regex    *string   `yaml:"regex"`
	In       *[]string `yaml:"in"`
	NotIn    *[]string `yaml:"not_in"`
}

var ErrRetrievingConfig = errors.New("error retrieving Config")

func GetConfig(configPath string) (Config, error) {
	var c *Config

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("%w: %v", ErrRetrievingConfig, err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return Config{}, fmt.Errorf("%w: %v", ErrRetrievingConfig, err)
	}

	return *c, nil
}
