package fs

import (
	"file-server/module/auth"
	"os"
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
		if vpath != string(os.PathSeparator) {
			vpathDepth = strings.Count(vpath, string(os.PathSeparator))
		}

		pathLinks = append(pathLinks, pathLink{
			vpath:      vpath,
			realpath:   realpath,
			vpathDepth: vpathDepth,
		})

		if vpath == string(os.PathSeparator) {
			continue
		}
		vpathParent := NormalizeVirtual(filepath.Dir(vpath))
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
			realpath = filepath.Join(pathLink.realpath, vpath[len(pathLink.vpath):])
			return realpath
		}
	}

	return _vpath
}

func (fs *Fs) Create(name string) (afero.File, error) {
	name = NormalizeVirtual(name)
	file, err := os.Create(fs.RealPath(name))
	if err != nil {
		return nil, err
	}
	return wrapFile(file, fs, name, IsRoot(name), fs.vpathChildren[name]), nil
}

func (fs *Fs) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(fs.RealPath(name), perm)
}

func (fs *Fs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(fs.RealPath(path), perm)
}

func (fs *Fs) Open(name string) (afero.File, error) {
	name = NormalizeVirtual(name)
	file, err := os.Open(fs.RealPath(name))
	if err != nil {
		return nil, err
	}
	return wrapFile(file, fs, name, IsRoot(name), fs.vpathChildren[name]), nil
}

func (fs *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	name = NormalizeVirtual(name)
	file, err := os.OpenFile(fs.RealPath(name), flag, perm)
	if err != nil {
		return nil, err
	}
	return wrapFile(file, fs, name, IsRoot(name), fs.vpathChildren[name]), nil
}

func (fs *Fs) Remove(name string) error {
	return os.Remove(fs.RealPath(name))
}

func (fs *Fs) RemoveAll(path string) error {
	return os.RemoveAll(fs.RealPath(path))
}

func (fs *Fs) Rename(oldname, newname string) error {
	return os.Rename(fs.RealPath(oldname), fs.RealPath(newname))
}

func (fs *Fs) Stat(name string) (os.FileInfo, error) {
	stat, err := os.Stat(fs.RealPath(name))
	if err != nil {
		return nil, err
	}
	return wrapStat(stat, filepath.Base(filepath.Clean(name))), nil
}

func (fs *Fs) Name() string {
	return "FileSystem"
}

func (fs *Fs) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(fs.RealPath(name), mode)
}

func (fs *Fs) Chown(name string, uid, gid int) error {
	return os.Chown(fs.RealPath(name), uid, gid)
}

func (fs *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(fs.RealPath(name), atime, mtime)
}

func (fs *Fs) ToWebdavFs() *WebdavFs {
	return &WebdavFs{
		Fs: fs,
	}
}
