package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liulihaocai/YetAnotherControlPanel/panel/i18n"
	"github.com/liulihaocai/YetAnotherControlPanel/util"
)

func RegisterWebPages(r *gin.Engine) error {
	err := InitializeLibraries(r)
	if err != nil {
		return err
	}
	err = InitializeThemes(r)
	if err != nil {
		return err
	}

	setupFS()

	log.Println("Registering web pages...")
	pageLogin(r)

	r.NoRoute(func(ctx *gin.Context) {
		// TODO: replace this with a real 404 page
		ctx.String(http.StatusNotFound, "404 Not Found")
	})

	return nil
}

func pageLogin(r *gin.Engine) {
	r.GET("/login", func(ctx *gin.Context) {
		tmpl := getBaseTemplate(ctx)
		tmpl.New("login").Parse(util.Must(util.ReadFile(templatesFS, "templates/login.html")))

		tmpl.Execute(ctx.Writer, struct {
			Basic TemplateInput
		}{
			Basic: TemplateInput{
				Title:  "Login",
				Locale: i18n.ReadLocale(ctx),
			},
		})
		ctx.Abort()
	})
}
