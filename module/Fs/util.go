package fs

import (
	"os"
	pathUtil "path"
	"path/filepath"
	"strings"
)

// vpath는 항상 '/'를 사용함
func NormalizeVirtual(path string) string {
	path = pathUtil.Clean(path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = pathUtil.Clean(path)
	return path
}

// realpath는 `os.PathSeparator`를 사용함
func NormalizeReal(path string) string {
	path = filepath.Clean(path)
	if !strings.HasPrefix(path, string(os.PathSeparator)) {
		cwd, err := os.Getwd()
		if err == nil {
			path = filepath.Join(cwd, path)
		} else {
			path = string(os.PathSeparator) + path
		}
	}
	path = filepath.Clean(path)
	return path
}

// 동일 경로도 부모로 칩니다.
func Isparent(parent, child string) bool {
	parent = filepath.Clean(parent)
	child = filepath.Clean(child)

	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}

	if strings.HasPrefix(rel, "..") {
		return false
	}

	return true
}

func IsRoot(path string) bool {
	path = filepath.Clean(path)
	return path == string(os.PathSeparator) || path == "/"
}
