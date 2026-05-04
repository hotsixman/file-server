package main

import (
	"file-server/module/LS"
	"file-server/module/auth"
	"file-server/module/ftp"
	"file-server/module/webdav"
	"log"
)

func main() {
	lockSystem := LS.NewMemLS()
	authManager, err := auth.LoadAuth("auth.json")
	if err != nil {
		log.Println(err)
		return
	}

	webdavApp := webdav.NewWebdavApp("test", lockSystem, authManager)
	ftpApp := ftp.NewFtpApp("test", lockSystem, authManager)
	go webdavApp.Listen("localhost:3000")
	go ftpApp.Listen("localhost:21")
	log.Println("Listen")

	select {}
}
