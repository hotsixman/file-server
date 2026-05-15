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

type NoPermissionError struct {
	Username string
	Path     string
	Action   string
}

func (e *NoPermissionError) Error() string {
	return fmt.Sprintf("permission denied: user %s, path %s, action %s", e.Username, e.Path, e.Action)
}
