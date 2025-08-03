package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func CreateFile(path, write string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, writeError := f.WriteString(write)

	if writeError != nil {
		return writeError
	}
	return nil
}

func AppendToFile(path, write string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, writeError := f.WriteString(write)
	if writeError != nil {
		return writeError
	}
	return nil
}

func FindInFile(path, search string) (bool, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, 0, err
	}
	defer f.Close()

	found := false
	repeated := 0
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, search) {
			found = true
			repeated++
		}
	}
	return found, repeated, scanner.Err()
}

func MakeDir(path, name string) error {
	err := os.MkdirAll(filepath.Join(path, name), 0755)
	if err != nil {
		return err
	}
	return nil
}

func IsFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsDirExist(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
