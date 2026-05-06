package fs

import (
	"fmt"
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

func (file *File) Readdir(count int) ([]os.FileInfo, error) {
	infos, err := file.File.Readdir(count)
	if err != nil {
		return nil, err
	}

	if file.virtualChildren != nil {
		for _, child := range file.virtualChildren {
			stat, err := file.fs.Stat(child)
			if err == nil {
				infos = append(infos, &FileInfo{
					FileInfo: stat,
					name:     filepath.Base(filepath.Clean(child)),
				})
			}
		}
	}

	fmt.Println(infos)
	if count <= 0 || len(infos) <= count {
		return infos, nil
	} else {
		return infos[:count], nil
	}
}

func (info *FileInfo) Name() string {
	return info.name
}
