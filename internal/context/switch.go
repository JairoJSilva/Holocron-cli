package context

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetStateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".flowctl")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "current"), nil
}

func SaveCurrent(client, cloud, account string) error {
	path, err := GetStateFilePath()
	if err != nil {
		return err
	}
	content := fmt.Sprintf("%s|%s|%s", client, cloud, account)
	return os.WriteFile(path, []byte(content), 0644)
}

func GetCurrent() (string, string, string, error) {
	path, err := GetStateFilePath()
	if err != nil {
		return "", "", "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", "", nil
		}
		return "", "", "", err
	}
	parts := strings.Split(string(data), "|")
	if len(parts) == 3 {
		return parts[0], parts[1], parts[2], nil
	}
	return "", "", "", nil
}

func ClearCurrent() error {
	path, err := GetStateFilePath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}
