package fs

import (
	"context"
	"os"

	"github.com/spf13/afero"
	"golang.org/x/net/webdav"
)

type WebdavFs struct {
	afero.Fs
}

func (wfs *WebdavFs) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return wfs.Fs.Mkdir(name, perm)
}

func (wfs *WebdavFs) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return wfs.Fs.OpenFile(name, flag, perm)
}

func (wfs *WebdavFs) RemoveAll(ctx context.Context, path string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return wfs.Fs.RemoveAll(path)
}

func (wfs *WebdavFs) Rename(ctx context.Context, oldname string, newname string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return wfs.Fs.Rename(oldname, newname)
}

func (wfs *WebdavFs) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return wfs.Fs.Stat(name)
}
