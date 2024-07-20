package file

import (
	"os"
	"strings"
)

var homeDir, _ = os.UserHomeDir()

func GetContents(dest string, file os.DirEntry) (string, error) {
	filePath := strings.Replace(dest, "~", homeDir, 1)

	openFile, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return "", err
	}

	defer func() { _ = openFile.Close() }()

	info, err := file.Info()
	if err != nil {
		return "", err
	}

	buf := make([]byte, info.Size())
	_, err = openFile.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func WriteTo(dest, code string) error {
	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer func() { _ = destFile.Close() }()

	_, err = destFile.WriteString(code)
	if err != nil {
		return err
	}

	return nil
}
