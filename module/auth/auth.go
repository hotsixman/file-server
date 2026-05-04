package auth

import (
	"encoding/json"
	"os"
	pathUtil "path"
	"slices"
	"strings"
)

type AuthUser struct {
	Password string   `json:"password"`
	Match    []string `json:"match"`
	Prefix   []string `json:"prefix"`
}

type Auth struct {
	json map[string]AuthUser
}

func LoadAuth(jsonPath string) (*Auth, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}

	auth := &Auth{
		json: map[string]AuthUser{},
	}
	err = json.Unmarshal(data, &auth.json)
	if err != nil {
		return nil, err
	}

	return auth, nil
}

func (auth *Auth) Authenticate(user, password string) bool {
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

func (auth *Auth) CheckPermission(user, path string) bool {
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
