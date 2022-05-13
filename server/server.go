package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liulihaocai/YetAnotherControlPanel/others"
)

var router *gin.Engine

func StartServer() {
	if others.TheConfig.Gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router = gin.Default()

	router.SetTrustedProxies(nil)

	if len(others.TheConfig.TrustedHosts) > 0 {
		router.Use(limitHosts(others.TheConfig.TrustedHosts))
		log.Println("Trusted hosts is enabled, only following hosts are allowed:", others.TheConfig.TrustedHosts)
	} else {
		log.Println("Trusted hosts is disabled, all hosts are allowed, this may be a security risk.")
	}

	api := router.Group("/api")
	setupApi(api)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, "404 Not Found")
	})

	go func() {
		err := router.Run(":" + strconv.Itoa(others.TheConfig.Port))
		if err != nil {
			panic(err)
		}
	}()

	log.Println("Server started on port", others.TheConfig.Port)
}

func setupApi(group *gin.RouterGroup) {

}
