package ftp

import (
	"errors"
	fs "file-server/module/Fs"
	"file-server/module/auth"
	"os"

	"github.com/spf13/afero"
)

type FtpClientDriver struct {
	afero.Fs
	user        string
	authManager auth.AuthManager
	rootDir     string
}

func newClientDriver(
	user string,
	rootDir string,
	authManager auth.AuthManager,
) *FtpClientDriver {
	driver := &FtpClientDriver{
		Fs: fs.NewFs(map[string]string{
			"/":        "test",
			"/foo":     "test2",
			"/foo/bar": "test3",
		}),
		user:        user,
		rootDir:     rootDir,
		authManager: authManager,
	}
	return driver
}

func (driver *FtpClientDriver) Create(path string) (afero.File, error) {
	if !driver.authManager.CheckPermission(driver.user, path) {
		return nil, errors.New("No Permission.")
	}
	return driver.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}

func (driver *FtpClientDriver) Open(path string) (afero.File, error) {
	if !driver.authManager.CheckPermission(driver.user, path) {
		return nil, errors.New("No Permission.")
	}
	return driver.OpenFile(path, os.O_RDONLY, 0)
}

func (driver *FtpClientDriver) OpenFile(path string, flag int, perm os.FileMode) (afero.File, error) {
	if !driver.authManager.CheckPermission(driver.user, path) {
		return nil, errors.New("No Permission.")
	}

	file, err := driver.Fs.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (driver *FtpClientDriver) Remove(path string) error {
	if !driver.authManager.CheckPermission(driver.user, path) {
		return errors.New("No Permission.")
	}

	return driver.Fs.Remove(path)
}

func (driver *FtpClientDriver) RemoveAll(path string) error {
	if !driver.authManager.CheckPermission(driver.user, path) {
		return errors.New("No Permission.")
	}

	return driver.Fs.RemoveAll(path)
}

func (driver *FtpClientDriver) Rename(oldpath, newpath string) error {
	if !driver.authManager.CheckPermission(driver.user, oldpath) || !driver.authManager.CheckPermission(driver.user, newpath) {
		return errors.New("No Permission.")
	}

	return driver.Fs.Rename(oldpath, newpath)
}
