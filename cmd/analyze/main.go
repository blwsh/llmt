package main

import (
	"context"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	envOpenAIToken = "OPENAI_TOKEN"
	envOllamaHost  = "OLLAMA_HOST"
)

var (
	cwd, _ = os.Getwd()
	ctx    = context.Background()
	log    = newLogger()
)

var (
	configPath        string
	defaultConfigPath = cwd + "/config.yaml"
)

func main() {
	rootCmd := &cobra.Command{Use: "llmt", Short: "llmt is a tool for analyzing projects"}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", defaultConfigPath, "config file path")

	rootCmd.AddCommand(&cobra.Command{
		Use:       "analyze",
		Aliases:   []string{"analyse", "a"},
		Short:     "analyze a project",
		Args:      cobra.ExactArgs(2),
		ValidArgs: []string{"source", "target"},
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := getConfig(configPath)
			if err != nil {
				log.Fatalf("failed to get config: %v", err)
			}

			source := strings.Replace(args[0], "~", os.Getenv("HOME"), 1)
			target := strings.Replace(args[1], "~", os.Getenv("HOME"), 1)

			if err := analyze(ctx, source, target, cfg); err != nil {
				log.Fatalf("failed to analyze project: %v", err)
			}
		},
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute command: %v", err)
	}
}
