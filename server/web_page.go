package server

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liulihaocai/YetAnotherControlPanel/others"
)

//go:embed all:nextjs/dist
var nextFS embed.FS

type WrappedFileSystem struct {
	fs     http.FileSystem
	config *others.Config
	prefix string
}

func (w *WrappedFileSystem) Open(name string) (http.File, error) {
	if name == "/"+w.config.SecuredEntrance {
		name = "/login"
	}
	name = w.prefix + name

	f, err := w.fs.Open(name)
	if err != nil {
		f, err1 := w.fs.Open(name + ".html")
		if err1 == nil {
			return f, nil
		}
		return nil, err
	}
	return f, nil
}

func setupWebPages(router *gin.Engine) {
	// TODO: i18n
	router.StaticFS("/", &WrappedFileSystem{
		fs:     http.FS(nextFS),
		config: others.TheConfig,
		prefix: "nextjs/dist",
	})
}
