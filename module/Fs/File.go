package fs

import (
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/afero"
)

type File struct {
	afero.File
	fs              *Fs
	vpath           string
	virtualChildren []string
	isRoot          bool
}

type FileInfo struct {
	os.FileInfo
	name string
}

func wrapFile(file afero.File, fs *Fs, vpath string, isRoot bool, virtualChildren []string) *File {
	return &File{File: file, fs: fs, vpath: vpath, isRoot: isRoot, virtualChildren: virtualChildren}
}

func wrapStat(stat os.FileInfo, name string) *FileInfo {
	return &FileInfo{
		FileInfo: stat,
		name:     name,
	}
}

func (file *File) Readdir(count int) ([]os.FileInfo, error) {
	infos, err := file.File.Readdir(count)
	if err != nil {
		return nil, err
	}

	if file.virtualChildren != nil {
		for _, child := range file.virtualChildren {
			stat, err := file.fs.Stat(child)
			if err == nil {
				infos = append(infos, wrapStat(stat, filepath.Base(filepath.Clean(child))))
			}
		}
	}

	if file.isRoot {
		readableDirs := file.fs.authManager.ReadableTopLevelDirs(file.fs.username)
		for i, v := range readableDirs {
			readableDirs[i] = filepath.Base(filepath.Clean(v))
		}
		infos = slices.DeleteFunc(infos, func(e os.FileInfo) bool {
			if !e.IsDir() {
				return false
			}
			return !slices.Contains(readableDirs, e.Name())
		})
	}

	if count <= 0 || len(infos) <= count {
		return infos, nil
	} else {
		return infos[:count], nil
	}
}

func (info *FileInfo) Name() string {
	return info.name
}
