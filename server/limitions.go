package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

func limitHosts(hosts []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(c.Request.Host)
		host := strings.Split(c.Request.Host, ":")[0]
		if !slices.Contains(hosts, host) {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}
