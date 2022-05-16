package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liulihaocai/YetAnotherControlPanel/others"
	"github.com/liulihaocai/YetAnotherControlPanel/util"
	"golang.org/x/exp/slices"
)

func limitHosts(hosts *[]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(*hosts) > 0 {
			host := c.Request.Host
			idx := strings.LastIndex(c.Request.Host, ":")
			if idx > 0 {
				host = c.Request.Host[:idx]
			}
			if !slices.Contains(*hosts, host) {
				c.AbortWithStatus(http.StatusForbidden)
			}
		}
	}
}

func limitLoginAndEntrance(conf *others.Config) gin.HandlerFunc {
	apiLogin := "api/login"
	allowedPageNotLogin := []*string{
		&conf.SecuredEntrance,
		&apiLogin,
	}

	return func(c *gin.Context) {
		session := ""
		// find session token in cookie
		cookie, err := c.Request.Cookie("session")
		if err == nil {
			session = cookie.Value
		}

		if !others.CheckSession(session) {
			if session != "" {
				// session token is invalid, remove it
				c.SetCookie("session", "", -1, "/", "", false, true)
			}

			if len(c.Request.URL.Path) < 1 || !util.MatchPointerSlice(allowedPageNotLogin, c.Request.URL.Path[1:]) {
				c.String(http.StatusUnauthorized, "401 Unauthorized, please login first.")
				c.Abort()
			}
		} else {
			if c.Request.URL.Path == "/" || c.Request.URL.Path == "//" {
				c.Redirect(http.StatusFound, "/index")
				c.Abort()
			}
		}
	}
}
