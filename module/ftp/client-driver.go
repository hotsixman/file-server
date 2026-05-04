package ftp

import (
	"errors"
	"file-server/module/auth"
	"os"
	"time"

	"github.com/spf13/afero"
	"golang.org/x/net/webdav"
)

type FtpClientDriver struct {
	afero.Fs
	lockSystem  webdav.LockSystem
	authManager *auth.Auth
	rootDir     string
	user        string
	password    string
}

func newClientDriver(rootDir string, lockSystem webdav.LockSystem, authManager *auth.Auth, user string, password string) *FtpClientDriver {
	driver := &FtpClientDriver{
		Fs:          afero.NewBasePathFs(afero.NewOsFs(), rootDir),
		rootDir:     rootDir,
		lockSystem:  lockSystem,
		authManager: authManager,
		user:        user,
		password:    password,
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

	if flag&os.O_WRONLY != 0 || flag&os.O_RDWR != 0 || flag&os.O_APPEND != 0 {
		release, err := driver.lockSystem.Confirm(time.Now(), path, "")
		if err != nil {
			return nil, err
		}
		release()
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

	release, err := driver.lockSystem.Confirm(time.Now(), path, "")
	if err != nil {
		return err
	}
	release()
	return driver.Fs.Remove(path)
}

func (driver *FtpClientDriver) RemoveAll(path string) error {
	if !driver.authManager.CheckPermission(driver.user, path) {
		return errors.New("No Permission.")
	}

	release, err := driver.lockSystem.Confirm(time.Now(), path, "")
	if err != nil {
		return err
	}
	release()
	return driver.Fs.RemoveAll(path)
}

func (driver *FtpClientDriver) Rename(oldpath, newpath string) error {
	if !driver.authManager.CheckPermission(driver.user, oldpath) || !driver.authManager.CheckPermission(driver.user, newpath) {
		return errors.New("No Permission.")
	}

	release, err := driver.lockSystem.Confirm(time.Now(), oldpath, newpath)
	if err != nil {
		return err
	}
	release()
	return driver.Fs.Rename(oldpath, newpath)
}
