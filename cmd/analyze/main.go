package main

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"cdb/pkg/project_analyzer"
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
				log.Fatal(err)
			}

			if err := analyze(ctx, args[0], args[1], cfg); err != nil {
				log.Fatal(err)
			}
		},
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func analyze(ctx context.Context, source, target string, c config) error {
	var (
		writer   = markdownWriter{logger: log}
		resolver = analyzerResolver{
			OpenAITokenResolver: func() string {
				var openAIToken string
				if openAIToken = os.Getenv("OPENAI_TOKEN"); os.Getenv("OPENAI_TOKEN") == "" {
					log.Fatal("OPENAI_TOKEN is not set")
				}

				return openAIToken
			},
			OllamaHostResolver: func() string {
				var ollamaHost string
				if ollamaHost = os.Getenv("OLLAMA_HOST"); os.Getenv("OLLAMA_HOST") == "" {
					ollamaHost = "http://localhost:11434"
				}

				return ollamaHost
			},
		}
		compiledRegexes = compileRegexesFromConfig(c)
		analyzers       []project_analyzer.FileAnalyzer
	)

	for _, a := range c.Analyzers {
		analyzer, err := resolver.resolve(a.Analyzer, a.Model)
		if err != nil {
			log.Warnf("failed to resolve analyzer %s: %v", a.Analyzer, err)
		}

		analyzers = append(analyzers, project_analyzer.FileAnalyzer{
			Prompt:        a.Prompt,
			Analyzer:      analyzer,
			Condition:     condition(a, compiledRegexes),
			ResultHandler: writer.write,
		})
	}

	err := project_analyzer.New(project_analyzer.WithLogger(log)).AnalyzeProject(ctx, source, target, analyzers)
	if err != nil {
		return err
	}

	return nil
}
