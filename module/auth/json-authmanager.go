package auth

import (
	"encoding/json"
	"os"
	pathUtil "path/filepath"
	"slices"
	"strings"
)

type AuthUser struct {
	Password string   `json:"password"`
	Match    []string `json:"match"`
	Prefix   []string `json:"prefix"`
}

type JsonAuth struct {
	json map[string]AuthUser
}

func LoadAuth(jsonPath string) (*JsonAuth, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}

	auth := &JsonAuth{
		json: map[string]AuthUser{},
	}
	err = json.Unmarshal(data, &auth.json)
	if err != nil {
		return nil, err
	}

	return auth, nil
}

func (auth *JsonAuth) Authenticate(user, password string) bool {
	authUser, ok := auth.json[user]
	if !ok {
		return false
	}

	if authUser.Password == password {
		return true
	} else {
		return false
	}
}

func (auth *JsonAuth) CheckPermission(user, path string) bool {
	authUser, ok := auth.json[user]
	if !ok {
		return false
	}

	path = pathUtil.Clean(path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	if slices.Contains(authUser.Match, path) {
		return true
	}
	for _, p := range authUser.Prefix {
		if !strings.HasPrefix(p, "/") {
			p = "/" + p
		}
		if !strings.HasSuffix(p, "/") {
			p += "/"
		}

		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
