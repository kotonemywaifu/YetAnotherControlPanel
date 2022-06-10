package util

import (
	"bytes"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
)

var minifier *minify.M

func InitializeMinifier() {
	minifier = minify.New()
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("text/javascript", js.Minify)
	minifier.AddFunc("application/json", json.Minify)
}

func MinifyCss(css string) (string, error) {
	return minifier.String("text/css", css)
}

func MinifyHtml(html string) (string, error) {
	return minifier.String("text/html", html)
}

func MinifyJs(js string) (string, error) {
	return minifier.String("text/javascript", js)
}

func StreamTemplateToUser(ctx *gin.Context, tmpl *template.Template, data any, useMinify bool) {
	if useMinify {
		var buf bytes.Buffer
		tmpl.Execute(&buf, data)
		res := buf.String()
		res, err := MinifyHtml(res)
		if err != nil {
			panic(err) // gin will handle this
		}
		ctx.Writer.Write([]byte(res))
	} else {
		tmpl.Execute(ctx.Writer, data)
	}
}
