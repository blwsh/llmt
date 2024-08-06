package main

import (
	"context"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/blwsh/llmt/cmd/llmt/config"
	"github.com/blwsh/llmt/cmd/llmt/handler/analyze"
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
		Use: "analyze", Aliases: []string{"analyse", "a"},
		Short: "analyze a project", Long: "analyze a project",
		Args: cobra.ExactArgs(2), ValidArgs: []string{"source", "target"},
		Run: func(cmd *cobra.Command, args []string) {
			fatalOnErr(analyze.Analyze(ctx, path(args[0]), path(args[1]), must(config.GetAnalyzeConfig(configPath))))
		},
	})

	fatalOnErr(rootCmd.Execute())
}

// path replaces ~ with the home directory
func path(path string) string {
	return strings.Replace(path, "~", os.Getenv("HOME"), 1)
}

// fatalOnErr logs the error and exits if the error is not nil
func fatalOnErr(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

// must take a tuple of (out, err) and return out if err is nil, otherwise log the error and exit
func must[T any](out T, err error) T {
	fatalOnErr(err)
	return out
}
