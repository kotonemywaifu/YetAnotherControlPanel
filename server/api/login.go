package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liulihaocai/YetAnotherControlPanel/others"
	"github.com/liulihaocai/YetAnotherControlPanel/panel/i18n"
	"github.com/liulihaocai/YetAnotherControlPanel/task"
)

func Login(group *gin.RouterGroup, cfg *others.Config) {
	group.POST("/login", func(c *gin.Context) {
		referer := c.Request.Referer()
		idx := strings.LastIndex(referer, "?")
		if idx > 0 {
			referer = referer[:idx]
		}
		if !strings.HasSuffix(referer, cfg.SecuredEntrance) {
			c.AbortWithStatus(http.StatusUnauthorized) // prevent behavior recognition
			return
		}
		// don't load locale if invalid secured entrance, used to get away from memory overflow attack
		locale := i18n.ReadLocale(c)

		if !ableToLogin(c.ClientIP()) {
			c.JSON(200, gin.H{
				"status": "error",
				"msg":    locale.Api.Login.FailedTooManyTimes,
			})
			return
		}

		accountHash := c.PostForm("account") // md5 encrypted (username+password)
		if len(accountHash) != 32 {
			c.JSON(200, gin.H{
				"status": "error",
				"msg":    locale.Api.Login.InvalidAccountHash,
			})
			failedIp(c.ClientIP())
			return
		}
		account := others.FindAccountHash(accountHash)
		if account == nil {
			c.JSON(200, gin.H{
				"status": "error",
				"msg":    locale.Api.Login.InvalidAccountCredentials,
			})
			failedIp(c.ClientIP())
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

	// add a task to process fail to ban
	task.AddTask(&task.Task{
		Name:      "panel-login-failToBan",
		NeedTicks: 5,
		Func:      failToBanTicker,
	})
}

var failures map[string]*LoginFailure = make(map[string]*LoginFailure)

func failToBanTicker() {
	now := time.Now()
	duration := time.Duration(others.TheConfig.FailToBan.BanTime) * time.Second
	for ip, failure := range failures {
		if now.Sub(failure.Time) > duration {
			failure.Time = now
			failure.Count--
		}
		if failure.Count <= 0 {
			delete(failures, ip)
		}
	}
}

func failedIp(ip string) {
	failure := failures[ip]
	if failure == nil {
		failure = &LoginFailure{
			Time:  time.Now(),
			Count: 1,
		}
		failures[ip] = failure
	} else {
		failure.Count++
	}
}

func ableToLogin(ip string) bool {
	failure := failures[ip]
	if failure == nil {
		return true
	}
	return failure.Count <= others.TheConfig.FailToBan.Failures
}

type LoginFailure struct {
	Count int
	Time  time.Time
}
