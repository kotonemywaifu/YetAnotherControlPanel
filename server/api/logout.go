package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/liulihaocai/YetAnotherControlPanel/others"
)

func Logout(group *gin.RouterGroup) {
	group.GET("/logout", func(c *gin.Context) {
		session, err := c.Cookie("session")
		if err != nil {
			session = ""
			log.Println("[LOGOUT] session cookie not found")
		}
		account := others.FindAccountSession(session)
		username := ""
		if account != nil {
			account.UpdateSession("EXPIRED", "")
			username = account.Username
		}
		c.SetCookie("session", account.Session, 3600, "/", "", false, true)
		c.JSON(200, gin.H{
			"status": "ok",
		})
		log.Println("Logout (User=" + username + ", IP=" + c.ClientIP() + ")")
	})
}
