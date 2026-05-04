package LS

import "golang.org/x/net/webdav"

func NewMemLS() webdav.LockSystem {
	return webdav.NewMemLS()
}
