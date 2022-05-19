package api

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liulihaocai/YetAnotherControlPanel/others"
)

func Login(group *gin.RouterGroup, cfg *others.Config) {
	group.POST("/login", func(c *gin.Context) {
		referer := c.Request.Referer()
		idx := strings.LastIndex(referer, "?")
		if idx > 0 {
			referer = referer[:idx]
		}
		if !strings.HasSuffix(referer, cfg.SecuredEntrance) {
			c.JSON(200, gin.H{
				"status": "error",
				"msg":    "invalid secured entrance",
			})
			return
		}

		accountHash := c.PostForm("account") // md5 encrypted (username+password)
		if len(accountHash) != 32 {
			c.JSON(200, gin.H{
				"status": "error",
				"msg":    "invalid account hash",
			})
			return
		}
		account := others.FindAccountHash(accountHash)
		if account == nil {
			c.JSON(200, gin.H{
				"status": "error",
				"msg":    "invalid username or password",
			})
			return
		}

		// correct password, set session cookie
		account.UpdateSession(c.Request.UserAgent())
		c.SetCookie("session", account.Session, 3600, "/", "", false, true)
		c.JSON(200, gin.H{
			"status": "ok",
		})
		log.Println("Login success (User=" + account.Username + ", IP=" + c.ClientIP() + ", UA=" + c.Request.UserAgent() + ")")
	})
}
