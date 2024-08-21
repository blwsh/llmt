package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AnalyzeConfig struct {
	Version   string                   `yaml:"version"`
	Analyzers []*ProjectAnalyzerConfig `yaml:"analyzers"`
}

type Assistant struct {
	Name        *string
	Description *string
	Files       *[]string
}

type ProjectAnalyzerConfig struct {
	Prompt    string     `yaml:"prompt"`
	Analyzer  string     `yaml:"analyzer"`
	Model     string     `yaml:"model"`
	Regex     *string    `yaml:"regex"`
	In        *[]string  `yaml:"in"`
	NotIn     *[]string  `yaml:"not_in"`
	Assistant *Assistant `yaml:"assistant"`
}

var ErrRetrievingConfig = errors.New("error retrieving Config")

func GetAnalyzeConfig(configPath string) (AnalyzeConfig, error) {
	var c *AnalyzeConfig

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return AnalyzeConfig{}, fmt.Errorf("%w: %v", ErrRetrievingConfig, err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return AnalyzeConfig{}, fmt.Errorf("%w: %v", ErrRetrievingConfig, err)
	}

	return *c, nil
}
