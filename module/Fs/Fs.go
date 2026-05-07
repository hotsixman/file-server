package fs

import (
	"file-server/module/auth"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/spf13/afero"
)

type Fs struct {
	afero.Fs
	pathLinks     []pathLink
	vpathChildren map[string][]string
	authManager   auth.AuthManager
	username      string
}

type pathLink struct {
	vpath      string
	realpath   string
	vpathDepth int
}

func NewFs(
	vpathMap map[string]string,
	authManager auth.AuthManager,
	username string,
) *Fs {
	pathLinks := make([]pathLink, 0)
	vpathChildren := map[string][]string{}

	for vpath, realpath := range vpathMap {
		vpath = NormalizeVirtual(vpath)
		realpath = NormalizeReal(realpath)

		vpathDepth := 0
		if !IsRoot(vpath) {
			vpathDepth = strings.Count(vpath, "/")
		}

		pathLinks = append(pathLinks, pathLink{
			vpath:      vpath,
			realpath:   realpath,
			vpathDepth: vpathDepth,
		})

		if IsRoot(vpath) {
			continue
		}
		vpathParent := NormalizeVirtual(filepath.ToSlash(filepath.Dir(vpath)))
		if vpath == vpathParent {
			continue
		}
		children, ok := vpathChildren[vpathParent]
		if !ok {
			children = make([]string, 0)
		}
		children = append(children, vpath)
		vpathChildren[vpathParent] = children
	}

	slices.SortFunc(pathLinks, func(a, b pathLink) int {
		return b.vpathDepth - a.vpathDepth
	})

	return &Fs{
		pathLinks:     pathLinks,
		vpathChildren: vpathChildren,
		authManager:   authManager,
		username:      username,
	}
}

func (fs *Fs) RealPath(vpath string) (realpath string) {
	_vpath := filepath.Clean(vpath)
	vpath = NormalizeVirtual(vpath)

	for _, pathLink := range fs.pathLinks {
		if Isparent(pathLink.vpath, vpath) {
			realpath = filepath.Clean(path.Join(pathLink.realpath, vpath[len(pathLink.vpath):]))
			return realpath
		}
	}

	return _vpath
}

func (fs *Fs) Create(vpath string) (afero.File, error) {
	vpath = path.Clean(vpath)
	if !path.IsAbs(vpath) {
		vpath = "/" + vpath
	}
	file, err := os.Create(fs.RealPath(vpath))
	if err != nil {
		return nil, err
	}

	return wrapFile(file, fs, vpath, IsRoot(vpath), fs.vpathChildren[vpath]), nil
}

func (fs *Fs) Mkdir(vpath string, perm os.FileMode) error {
	return os.Mkdir(fs.RealPath(vpath), perm)
}

func (fs *Fs) MkdirAll(vpath string, perm os.FileMode) error {
	return os.MkdirAll(fs.RealPath(vpath), perm)
}

func (fs *Fs) Open(vpath string) (afero.File, error) {
	vpath = path.Clean(vpath)
	if !path.IsAbs(vpath) {
		vpath = "/" + vpath
	}
	file, err := os.Open(fs.RealPath(vpath))
	if err != nil {
		return nil, err
	}
	return wrapFile(file, fs, vpath, IsRoot(vpath), fs.vpathChildren[vpath]), nil
}

func (fs *Fs) OpenFile(vpath string, flag int, perm os.FileMode) (afero.File, error) {
	vpath = path.Clean(vpath)
	if !path.IsAbs(vpath) {
		vpath = "/" + vpath
	}
	file, err := os.OpenFile(fs.RealPath(vpath), flag, perm)
	if err != nil {
		return nil, err
	}
	return wrapFile(file, fs, vpath, IsRoot(vpath), fs.vpathChildren[vpath]), nil
}

func (fs *Fs) Remove(vpath string) error {
	vpath = path.Clean(vpath)
	if !path.IsAbs(vpath) {
		vpath = "/" + vpath
	}
	return os.Remove(fs.RealPath(vpath))
}

func (fs *Fs) RemoveAll(vpath string) error {
	vpath = path.Clean(vpath)
	if !path.IsAbs(vpath) {
		vpath = "/" + vpath
	}
	return os.RemoveAll(fs.RealPath(vpath))
}

func (fs *Fs) Rename(oldvpath, newvpath string) error {
	oldvpath = path.Clean(oldvpath)
	if !path.IsAbs(oldvpath) {
		oldvpath = "/" + oldvpath
	}
	newvpath = path.Clean(newvpath)
	if !path.IsAbs(newvpath) {
		newvpath = "/" + newvpath
	}
	return os.Rename(fs.RealPath(oldvpath), fs.RealPath(newvpath))
}

func (fs *Fs) Stat(vpath string) (os.FileInfo, error) {
	vpath = path.Clean(vpath)
	if !path.IsAbs(vpath) {
		vpath = "/" + vpath
	}
	stat, err := os.Stat(fs.RealPath(vpath))
	if err != nil {
		return nil, err
	}
	return wrapStat(stat, filepath.Base(filepath.Clean(vpath))), nil
}

func (fs *Fs) Name() string {
	return "FileSystem"
}

func (fs *Fs) Chmod(vpath string, mode os.FileMode) error {
	vpath = path.Clean(vpath)
	if !path.IsAbs(vpath) {
		vpath = "/" + vpath
	}
	return os.Chmod(fs.RealPath(vpath), mode)
}

func (fs *Fs) Chown(vpath string, uid, gid int) error {
	vpath = path.Clean(vpath)
	if !path.IsAbs(vpath) {
		vpath = "/" + vpath
	}
	return os.Chown(fs.RealPath(vpath), uid, gid)
}

func (fs *Fs) Chtimes(vpath string, atime time.Time, mtime time.Time) error {
	vpath = path.Clean(vpath)
	if !path.IsAbs(vpath) {
		vpath = "/" + vpath
	}
	return os.Chtimes(fs.RealPath(vpath), atime, mtime)
}

func (fs *Fs) ToWebdavFs() *WebdavFs {
	return &WebdavFs{
		Fs: fs,
	}
}
