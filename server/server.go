package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liulihaocai/YetAnotherControlPanel/others"
	"github.com/liulihaocai/YetAnotherControlPanel/server/api"
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
		log.Println("Trusted hosts is enabled, only following hosts are allowed:", others.TheConfig.TrustedHosts)
	} else {
		log.Println("Trusted hosts is disabled, all hosts are allowed, this may be a security risk.")
	}
	router.Use(limitHosts(&others.TheConfig.TrustedHosts))

	if others.TheConfig.SecuredEntrance == "" {
		log.Println("Secured entrance is disabled, this may be a security risk.")
	} else {
		log.Println("Secured entrance is enabled, you can only login from the secured entrance: /", others.TheConfig.SecuredEntrance)
	}
	router.Use(limitLoginAndEntrance(others.TheConfig))

	api := router.Group("/api")
	setupApi(api, others.TheConfig)

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

func setupApi(group *gin.RouterGroup, cfg *others.Config) {
	api.Login(group, cfg)
}
