package main

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	Version   string                   `yaml:"version"`
	Analyzers []*projectAnalyzerConfig `yaml:"analyzers"`
}

type projectAnalyzerConfig struct {
	Prompt   string    `yaml:"prompt"`
	Analyzer string    `yaml:"analyzer"`
	Model    string    `yaml:"model"`
	Regex    *string   `yaml:"regex"`
	In       *[]string `yaml:"in"`
	NotIn    *[]string `yaml:"not_in"`
}

var ErrRetrievingConfig = errors.New("error retrieving config")

func getConfig(configPath string) (config, error) {
	var c *config

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return config{}, fmt.Errorf("%w: %v", ErrRetrievingConfig, err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return config{}, fmt.Errorf("%w: %v", ErrRetrievingConfig, err)
	}

	return *c, nil
}
