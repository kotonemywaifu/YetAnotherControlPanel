package api

import (
	"github.com/gin-gonic/gin"
	"github.com/liulihaocai/YetAnotherControlPanel/others"
	"github.com/liulihaocai/YetAnotherControlPanel/util"
)

func Login(group *gin.RouterGroup, cfg *others.Config) {
	group.POST("/login", func(c *gin.Context) {
		entrance := c.PostForm("entrance")
		if entrance != cfg.SecuredEntrance {
			c.JSON(200, gin.H{
				"status": "error",
				"msg":    "invalid secured entrance",
			})
			return
		}

		passwd := c.PostForm("password") // md5 encrypted password
		if len(passwd) != 32 || passwd != util.GetMD5Hash(cfg.Password) {
			c.JSON(200, gin.H{
				"status": "error",
				"msg":    "invalid password",
			})
			return
		}

		// correct password, set session cookie
		c.SetCookie("session", cfg.Session, 3600, "/", "", false, true)
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
