package auth

import "slices"

func (am *AuthManager) Authenticate(username, password string) bool {
	var count int
	err := am.db.QueryRow("SELECT COUNT(*) AS COUNT FROM User Where username = ? AND password = ?", username, password).Scan(&count)
	if err != nil {
		return false
	}

	if count > 0 {
		return true
	} else {
		return false
	}
}

func (am *AuthManager) Permission(username, path string) Permission
func (am *AuthManager) Permissions(username, paths []string) map[string]Permission
func (am *AuthManager) ReadPermission(username, path string) bool
func (am *AuthManager) ReadPermissions(username, paths []string) map[string]bool
func (am *AuthManager) WritePermission(username, path string) bool
func (am *AuthManager) WritePermissions(username, paths []string) map[string]bool

func (am *AuthManager) DirMap(username string, dirMap map[string]string) map[string]string {
	readableDirs := am.ReadableTopLevelDirs(username)
	readableDirMap := map[string]string{}

	for v, r := range dirMap {
		if v == "/" {
			readableDirMap[v] = r
		} else if slices.Contains(readableDirs, v) {
			readableDirMap[v] = r
		}
	}

	return readableDirMap
}

func (am *AuthManager) ReadableTopLevelDirs(username string) []string {
	userId, err := am.getUserId(username)
	if err != nil {
		return []string{}
	}
	groupIds := am.getUserGroupIds(userId)
	dirs := am.getReadableTopLevelDirs(userId, groupIds)
	return dirs
}
