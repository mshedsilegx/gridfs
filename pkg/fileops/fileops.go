package fileops

import (
	"bufio"
	"os"
	"strings"
)

// FileExistsAndNotEmpty checks if the file exists and is not empty.
func FileExistsAndNotEmpty(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) || info.Size() == 0 {
		return false
	}
	return true
}

// ReadFileNames reads file names from the given file.
func ReadFileNames(filename string) (names []string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			names = append(names, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return names, nil
}

// CreateDirectory creates a directory if it doesn't exist.
func CreateDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}
