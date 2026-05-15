package ftp

import (
	"file-server/module/auth"
	"net"

	ftpserver "github.com/fclairamb/ftpserverlib"
)

type FtpApp struct {
	server      *ftpserver.FtpServer
	driver      *FtpMainDriver
	authManager *auth.AuthManager
}

func NewFtpApp(dirMap map[string]string, authManager *auth.AuthManager) *FtpApp {
	driver := newMainDriver(dirMap)
	app := &FtpApp{
		server:      ftpserver.NewFtpServer(driver),
		driver:      driver,
		authManager: authManager,
	}
	return app
}

func (app *FtpApp) Listen(address string) error {
	app.driver.address = address
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	app.driver.listener = listener
	app.driver.authManager = app.authManager
	return app.server.ListenAndServe()
}
