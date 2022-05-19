package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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
		log.Println("Secured entrance is enabled, you can only login from the secured entrance: /" + others.TheConfig.SecuredEntrance)
	}
	router.Use(limitLoginAndEntrance(others.TheConfig))

	setupWebPages(router)
	api := router.Group("/api")
	setupApi(api, others.TheConfig)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, "404 Not Found")
	})

	serveByRouter(router)
}

func serveByRouter(router *gin.Engine) {
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(others.TheConfig.Port),
		Handler: router,
	}

	go func() {
		var err error
		if others.TheConfig.Gin.Tls.Enabled {
			srv.ListenAndServeTLS(others.TheConfig.Gin.Tls.CertFile, others.TheConfig.Gin.Tls.KeyFile)
		} else {
			log.Println("You are running in non-secure mode, please consider using TLS.")
			err = srv.ListenAndServe()
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panicln(err)
		}
	}()

	log.Println("Server started on port", others.TheConfig.Port)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("\rShutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	os.Exit(0)
}

func setupApi(group *gin.RouterGroup, cfg *others.Config) {
	api.Login(group, cfg)
}
