package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterWebPages(r *gin.Engine) error {
	err := InitializeLibraries()
	if err != nil {
		return err
	}

	log.Println("Registering web pages...")

	r.NoRoute(func(ctx *gin.Context) {
		// TODO: replace this with a real 404 page
		ctx.String(http.StatusNotFound, "404 Not Found")
	})

	return nil
}
