package auth

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type AuthManager struct {
	db *sql.DB
}

func NewAuthManager() (*AuthManager, error) {
	db, err := sql.Open("sqlite", "main.db")
	if err != nil {
		return nil, err
	}

	am := &AuthManager{
		db: db,
	}
	err = am.setup()

	if err != nil {
		return nil, err
	}

	return am, nil
}

//setup

func (am *AuthManager) setup() error {
	tx, err := am.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS User (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	username TEXT NOT NULL UNIQUE,
    	password TEXT NOT NULL
	);`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS "Group" (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	name TEXT NOT NULL UNIQUE
	);`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS GroupMember (
		groupId INTEGER NOT NULL,
		userId INTEGER NOT NULL,
		PRIMARY KEY (groupId, userId),
		FOREIGN KEY (groupId) REFERENCES "Group"(id) ON DELETE CASCADE,
		FOREIGN KEY (userId) REFERENCES User(id) ON DELETE CASCADE
	);`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS UserPermission (
	    userId INTEGER NOT NULL,
	    dirname TEXT NOT NULL,
	    read BOOLEAN NOT NULL CHECK (read IN (0, 1)),
	    write BOOLEAN NOT NULL CHECK (write IN (0, 1)),
	    PRIMARY KEY (userId, dirname),
	    FOREIGN KEY (userId) REFERENCES User(id) ON DELETE CASCADE
	);`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS GroupPermission (
	    groupId INTEGER NOT NULL,
	    dirname TEXT NOT NULL,
	    read BOOLEAN NOT NULL CHECK (read IN (0, 1)),
	    write BOOLEAN NOT NULL CHECK (write IN (0, 1)),
	    PRIMARY KEY (groupId, dirname),
	    FOREIGN KEY (groupId) REFERENCES "Group"(id) ON DELETE CASCADE
	);`)
	if err != nil {
		return err
	}

	userId, err := am.setupAdminUser(tx)
	if err != nil {
		return err
	}
	err = am.setupAdminGroup(tx, userId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

// 어드민 유저를 새로 만들었을 때에만 userId가 -1이 아님
func (am *AuthManager) setupAdminUser(tx *sql.Tx) (userId int64, err error) {
	var count int

	err = tx.QueryRow("SELECT COUNT(*) FROM User").Scan(&count)
	if err != nil {
		return -1, err
	}

	if count > 0 {
		return -1, nil
	}

	result, err := tx.Exec("INSERT INTO `User` (`username`, `password`) VALUES ('admin', '0000')")
	if err != nil {
		return -1, nil
	}
	userId, err = result.LastInsertId()
	if err != nil {
		return -1, nil
	}

	return userId, nil
}

func (am *AuthManager) setupAdminGroup(tx *sql.Tx, userId int64) error {
	var groupId int64
	err := tx.QueryRow("SELECT `id` FROM `Group` WHERE `name` = 'admin'").Scan(&groupId)
	if err != nil {
		if err == sql.ErrNoRows {
			result, err := tx.Exec("INSERT INTO `Group` (`name`) VALUES ('admin')")
			if err != nil {
				return err
			}
			groupId, err = result.LastInsertId()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if userId != -1 {
		_, err = tx.Exec("INSERT INTO `GroupMember` (`groupId`, `userId`) VALUES (?, ?)", groupId, userId)
		if err != nil {
			return err
		}
	}

	return nil
}

// DB 형식

type User struct {
	id       int
	username string
	password string
}
type Group struct {
	id   int
	name string
}
type GroupMember struct {
	groupId int
	userId  int
}
type UserPermission struct {
	userId int
	// `/`로 시작해야함
	dirname string
	read    bool
	write   bool
}
type GroupPermission struct {
	groupId int
	// `/`로 시작해야함
	dirname string
	read    bool
	write   bool
}
