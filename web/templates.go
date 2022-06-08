package web

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/liulihaocai/YetAnotherControlPanel/others"
	"github.com/liulihaocai/YetAnotherControlPanel/panel/i18n"
	"github.com/liulihaocai/YetAnotherControlPanel/util"
)

//go:embed templates/*.html
var embedFS embed.FS

var templatesFS fs.FS

type TemplateInput struct {
	Title  string
	Locale *i18n.Locale
}

func getBaseTemplate(ctx *gin.Context) *template.Template {
	tmpl := template.Must(template.New("base").Parse(util.Must(util.ReadFile(templatesFS, "templates/base.html"))))

	scripts := tryPushLibrary(ctx)
	if scripts != "" {
		tmpl.New("scripts").Parse(scripts)
	}

	return tmpl
}

func setupFS() {
	if others.TheConfig.Gin.UseDebugFS {
		log.Println("Using debug filesystem...")
		templatesFS = os.DirFS("./web")
	} else {
		templatesFS = embedFS
	}
}
