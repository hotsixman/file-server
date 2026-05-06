package fs

import (
	"os"
	"path/filepath"
	"strings"
)

func normalizeVirtual(path string) string {
	path = filepath.Clean(path)
	if !strings.HasPrefix(path, string(os.PathSeparator)) {
		path = string(os.PathSeparator) + path
	}
	if !strings.HasSuffix(path, string(os.PathSeparator)) {
		path += string(os.PathSeparator)
	}
	return path
}
func normalizeReal(path string) string {
	path = filepath.Clean(path)
	if !strings.HasPrefix(path, string(os.PathSeparator)) {
		cwd, err := os.Getwd()
		if err == nil {
			path = filepath.Join(cwd, path)
		} else {
			path = string(os.PathSeparator) + path
		}
	}
	if !strings.HasSuffix(path, string(os.PathSeparator)) {
		path += string(os.PathSeparator)
	}
	return path
}
