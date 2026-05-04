package main

import (
	"file-server/module/LS"
	"file-server/module/auth"
	"file-server/module/ftp"
	"file-server/module/webdav"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	lockSystem := LS.NewMemLS()
	authManager, err := auth.LoadAuth("auth.json")
	if err != nil {
		log.Println(err)
		return
	}

	webdavApp := webdav.NewWebdavApp(os.Getenv("ROOT_DIR"), lockSystem, authManager)
	ftpApp := ftp.NewFtpApp(os.Getenv("ROOT_DIR"), lockSystem, authManager)
	go webdavApp.Listen("localhost:3000")
	go ftpApp.Listen("localhost:21")
	log.Println("Listen")

	select {}
}
