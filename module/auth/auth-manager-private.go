package auth

func (am *AuthManager) getUserId(username string) (int, error) {
	var id int
	err := am.db.QueryRow("SELECT `id` FROM `User` WHERE `username` = ?", username).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

/*
유저가 속한 그룹 Id를 모두 가져옴
*/
func (am *AuthManager) getUserGroupIds(userId int) []int {
	rows, err := am.db.Query("SELECT `groupId` FROM `GroupMember` WHERE `userId` = ?", userId)
	groupIds := make([]int, 0)
	if err != nil {
		return groupIds
	}

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			continue
		}
		groupIds = append(groupIds, id)
	}

	return groupIds
}

/*
허용된 top level 디렉토리를 가져옴
*/
func (am *AuthManager) getReadableTopLevelDirs(userId int, groupIds []int) []string {
	setToSlice := func(set map[string]bool) []string {
		s := make([]string, 0)
		for k := range set {
			s = append(s, k)
		}
		return s
	}

	dirSet := map[string]bool{}

	rows, err := am.db.Query("SELECT `dirname` FROM `UserPermission` WHERE `userId` = ? AND `read` = 1", userId)
	if err != nil {
		return setToSlice(dirSet)
	}

	for rows.Next() {
		var dirname string
		err := rows.Scan(&dirname)
		if err != nil {
			continue
		}
		dirSet[dirname] = true
	}

	for _, groupId := range groupIds {
		rows, err := am.db.Query("SELECT `dirname` FROM `GroupPermission` WHERE `groupId` = ? AND `read` = 1", groupId)
		if err != nil {
			return setToSlice(dirSet)
		}

		for rows.Next() {
			var dirname string
			err := rows.Scan(&dirname)
			if err != nil {
				continue
			}
			dirSet[dirname] = true
		}
	}

	return setToSlice(dirSet)
}
