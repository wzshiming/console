//go:generate go-bindata -o static_binddata.go --pkg static --prefix web web/...

package static

import (
	"net/http"

	"github.com/wzshiming/go-bindata/fs"
)

func NewFileSystem() http.FileSystem {
	return &fs.AssetFS{
		Asset: Asset,
		Index: "index.html",
	}
}
