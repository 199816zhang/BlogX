package log_service

import (
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"fmt"
	"github.com/gin-gonic/gin"
)

func NewLoginSuccess(c *gin.Context, loginType enum.LoginType) {
	ip := c.ClientIP()
	addr := core.GetIpAddr(ip)
	token := c.GetHeader("token")
	fmt.Println(token)
	UserID := 1
	UserName := ""
	global.DB.Create(&models.LogModel{
		LogType:   enum.LoginLogType,
		Title:     "登录成功",
		Content:   "",
		UserID:    uint(UserID),
		UserModel: models.UserModel{},
		IP:        ip,
		Addr:      addr,
		LogStatus: true,
		UserName:  UserName,
		Pwd:       "-",
		LoginType: loginType,
	})
}
func NewLoginFail(c *gin.Context, loginType enum.LoginType, msg string, username string, pwd string) {
	ip := c.ClientIP()
	addr := core.GetIpAddr(ip)
	global.DB.Create(&models.LogModel{
		LogType:   enum.LoginLogType,
		Title:     "登录失败",
		Content:   msg,
		UserModel: models.UserModel{},
		IP:        ip,
		Addr:      addr,
		LogStatus: false,
		UserName:  username,
		Pwd:       pwd,
		LoginType: loginType,
	})
}
