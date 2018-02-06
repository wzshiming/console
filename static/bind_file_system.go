//go:generate go-bindata -o static_binddata.go -pkg static -ignore '^\.' ./web/...

package static

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

type File struct {
	io.Reader
	io.Seeker
	ReaddirFunc func(count int) ([]os.FileInfo, error)
	StatFunc    func() (os.FileInfo, error)
}

func (f File) Close() error {
	return nil
}

func (f File) Readdir(count int) ([]os.FileInfo, error) {
	return f.ReaddirFunc(count)
}

func (f File) Stat() (os.FileInfo, error) {
	return f.StatFunc()
}

type FileSystem struct {
	Dir  string
	Root string
}

func NewFileSystem() *FileSystem {
	return &FileSystem{
		Dir:  "web",
		Root: "index.html",
	}
}

func (f *FileSystem) Open(name string) (http.File, error) {
	if strings.HasSuffix(name, "/") {
		name = path.Join(f.Dir, name, f.Root)
	} else {
		name = path.Join(f.Dir, name)
	}

	d, err := Asset(name)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(d)

	return File{
		reader,
		reader,
		func(count int) ([]os.FileInfo, error) {
			fi := []os.FileInfo{}
			ps, err := AssetDir(path.Dir(name))
			if err != nil {
				return nil, err
			}
			for i := 0; i != count; i++ {
				ai, err := AssetInfo(ps[i])
				if err != nil {
					return nil, err
				}
				fi = append(fi, ai)
			}
			return fi, nil
		},
		func() (os.FileInfo, error) {
			return AssetInfo(name)
		},
	}, nil
}
