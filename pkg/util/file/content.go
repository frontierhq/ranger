package file

import (
	"os"
	"path/filepath"
)

func CreateOrUpdate(fullPath string, content string, append bool) error {
	dir := filepath.Dir(fullPath)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	flags := os.O_CREATE | os.O_WRONLY
	if append {
		flags = flags | os.O_APPEND
	}

	file, err := os.OpenFile(fullPath, flags, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func Clear(fullPath string) error {
	err := os.Remove(fullPath)
	if err != nil {
		return err
	}

	return nil
}
