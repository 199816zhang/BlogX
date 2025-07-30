package site_api

import (
	"blogx_server/models/enum"
	"blogx_server/service/log_service"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SiteApi struct{}

func (SiteApi) SiteInfoView(c *gin.Context) {
	fmt.Println("SiteInfoView")
	log_service.NewLoginSuccess(c, enum.UserPwdLoginType)
	log_service.NewLoginFail(c, enum.UserPwdLoginType, "用户不存在", "zh", "123456")
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": nil,
	})
	return
}
