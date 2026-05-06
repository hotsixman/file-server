package ftp

import (
	"crypto/tls"
	"errors"
	"file-server/module/auth"
	"net"

	ftpserverlib "github.com/fclairamb/ftpserverlib"
)

type FtpMainDriver struct {
	ftpserverlib.MainDriver
	rootDir     string
	authManager auth.AuthManager
	address     string
	listener    net.Listener
}

func newMainDriver(rootDir string) *FtpMainDriver {
	return &FtpMainDriver{rootDir: rootDir}
}

func (driver *FtpMainDriver) GetSettings() (*ftpserverlib.Settings, error) {
	settings := &ftpserverlib.Settings{
		ListenAddr: driver.address,
		Listener:   driver.listener,
	}
	return settings, nil
}

func (driver *FtpMainDriver) ClientConnected(cc ftpserverlib.ClientContext) (string, error) {
	return "", nil
}

func (driver *FtpMainDriver) ClientDisconnected(cc ftpserverlib.ClientContext) {}

func (driver *FtpMainDriver) AuthUser(cc ftpserverlib.ClientContext, user, pass string) (ftpserverlib.ClientDriver, error) {
	if !driver.authManager.Authenticate(user, pass) {
		return nil, errors.New("Invalid User or Password.")
	}
	return newClientDriver(user, driver.rootDir, driver.authManager), nil
}

func (driver *FtpMainDriver) GetTLSConfig() (*tls.Config, error) {
	return nil, errors.New("No TLS")
}
