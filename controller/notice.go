package controller

import (
	"github.com/gin-gonic/gin"
	"goflylivechat/common"
	"goflylivechat/models"
	"strings"
)

func GetNotice(c *gin.Context) {
	kefuId := c.Query("kefu_id")
	user := models.FindUser(kefuId)
	if user.ID == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "user not found",
		})
		return
	}
	welcomeMessage := models.FindConfigByUserId(user.Name, "WelcomeMessage")
	offlineMessage := models.FindConfigByUserId(user.Name, "OfflineMessage")
	allNotice := models.FindConfigByUserId(user.Name, "AllNotice")
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"welcome":   visitorNoticeText("WelcomeMessage", welcomeMessage.ConfValue),
			"offline":   visitorNoticeText("OfflineMessage", offlineMessage.ConfValue),
			"avatar":    common.KefuAvatar,
			"nickname":  common.PublicKefuName,
			"allNotice": visitorNoticeText("AllNotice", allNotice.ConfValue),
		},
	})
}

func visitorNoticeText(key, value string) string {
	value = strings.TrimSpace(value)
	switch key {
	case "WelcomeMessage":
		if value == "" || value == "How may I help you?" {
			return common.DefaultWelcome
		}
	case "OfflineMessage":
		if value == "" || value == "I am currently offline and will reply to you later!" {
			return common.DefaultOffline
		}
	case "AllNotice":
		if value == "" || value == "Open source customer support system at your service" {
			return common.DefaultAllNotice
		}
	}
	return value
}
