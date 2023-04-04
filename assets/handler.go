package assets

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"
)

const defaultUrl = "index.html"

type pseudoFs func(name string) (fs.File, error)

func (f pseudoFs) Open(name string) (fs.File, error) {
	return f(name)
}

func Handler(resources embed.FS, urlPrefix string, resourcesRoot string) http.Handler {
	handler := pseudoFs(func(name string) (fs.File, error) {
		assetPath := path.Join(resourcesRoot, name)

		f, err := resources.Open(assetPath)
		if os.IsNotExist(err) {
			return resources.Open(defaultUrl)
		}

		return f, err
	})

	return http.StripPrefix(urlPrefix, http.FileServer(http.FS(handler)))
}
