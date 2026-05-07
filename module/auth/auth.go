package auth

import "fmt"

type Permission struct {
	Read  bool
	Write bool
}

/*
모든 경로는 '/'를 사용합니다.
filepath 패키지 대신 path 패키지를 사용하십시오.
*/
type AuthManager interface {
	Authenticate(username, password string) bool

	Permission(username, path string) Permission
	Permissions(username, paths []string) map[string]Permission
	ReadPermission(username, path string) bool
	ReadPermissions(username, paths []string) map[string]bool
	WritePermission(username, path string) bool
	WritePermissions(username, paths []string) map[string]bool

	DirMap(rootDir string, username string) map[string]string
	AllowedDirs(username string) []string
}

type NoPermissionError struct {
	Username string
	Path     string
	Action   string
}

func (e *NoPermissionError) Error() string {
	return fmt.Sprintf("permission denied: user %s, path %s, action %s", e.Username, e.Path, e.Action)
}
