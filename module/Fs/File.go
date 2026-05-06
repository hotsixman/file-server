package fs

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type File struct {
	afero.File
	fs              *Fs
	virtualChildren []string
}

type FileInfo struct {
	os.FileInfo
	name string
}

func wrapFile(file afero.File, fs *Fs, virtualChildren []string) *File {
	return &File{File: file, fs: fs, virtualChildren: virtualChildren}
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

	if count <= 0 || len(infos) <= count {
		return infos, nil
	} else {
		return infos[:count], nil
	}
}

func (info *FileInfo) Name() string {
	return info.name
}
