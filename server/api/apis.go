package api

import (
	"github.com/gin-gonic/gin"
	"github.com/liulihaocai/YetAnotherControlPanel/others"
)

func RegisterApi(group *gin.RouterGroup, cfg *others.Config) {
	Login(group, cfg)
}
