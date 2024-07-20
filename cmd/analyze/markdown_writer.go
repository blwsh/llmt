package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blwsh/llmt/lib/file"
	"github.com/blwsh/llmt/lib/logger"
)

type markdownWriter struct {
	logger logger.Logger
}

func (m *markdownWriter) write(destFilepath string, result string) error {
	if _, err := os.Stat(filepath.Dir(destFilepath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(destFilepath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	outputPath := filepath.Join(destFilepath + ".md")

	err := file.WriteTo(outputPath, result)
	if err != nil {
		m.logger.Error(err)
	}

	m.logger.Info("analyzed file: ", filepath.Base(destFilepath), " -> ", outputPath)

	return nil
}
