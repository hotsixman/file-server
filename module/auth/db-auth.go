package auth

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type SqliteAM struct {
	AuthManager
	db *sql.DB
}

func NewSqliteAM() (*SqliteAM, error) {
	db, err := sql.Open("sqlite", "main.db")
	if err != nil {
		return nil, err
	}

	am := &SqliteAM{
		db: db,
	}
	err = am.setup()

	if err != nil {
		return nil, err
	}

	return am, nil
}

func (am *SqliteAM) Authenticate(username, password string) bool {
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

func (am *SqliteAM) Permission(username, path string) bool {
	return true
}

func (am *SqliteAM) DirMap(rootDir string, username string) map[string]string {
	return map[string]string{
		"/":    rootDir,
		"/foo": "test2",
	}
}

func (am *SqliteAM) AllowedTopLevelDirs(username string) []string {
	return []string{"반", "foo"}
}

func (am *SqliteAM) setup() error {
	_, err := am.db.Exec(`
	CREATE TABLE IF NOT EXISTS User (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	username TEXT NOT NULL UNIQUE,
    	password TEXT NOT NULL
	);`)
	if err != nil {
		return err
	}
	_, err = am.db.Exec(`
	CREATE TABLE IF NOT EXISTS "Group" (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	name TEXT NOT NULL UNIQUE
	);`)
	if err != nil {
		return err
	}
	_, err = am.db.Exec(`
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
	_, err = am.db.Exec(`
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
	_, err = am.db.Exec(`
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

	var count int
	err = am.db.QueryRow("SELECT COUNT(*) FROM User").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = am.db.Exec("INSERT INTO `User` (`username`, `password`) VALUES ('admin', '1234')")
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
	userId  int
	dirname string
	read    bool
	write   bool
}
type GroupPermission struct {
	groupId int
	dirname string
	read    bool
	write   bool
}
