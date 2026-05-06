package auth

type AuthManager interface {
	Authenticate(username, password string) bool
	CheckPermission(username, path string) bool
}
