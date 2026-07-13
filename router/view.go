package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitViewRouter(engine *gin.Engine) {
	//engine.GET("/", tmpl.PageIndex)
	engine.GET("/login", PageLogin)
	engine.GET("/pannel", PagePannel)
	engine.GET("/livechat", PageChat)
	engine.GET("/main", PageMain)
	engine.GET("/chat_main", PageChatMain)
	engine.GET("/m/chat", PageMobileChat)
	engine.GET("/chat_mobile", PageMobileChat)
	engine.GET("/setting", PageSetting)
}

// Login page
func PageLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// Dashboard
func PagePannel(c *gin.Context) {
	c.HTML(http.StatusOK, "pannel.html", nil)
}

// Admin console
func PageMain(c *gin.Context) {
	c.HTML(http.StatusOK, "main.html", nil)
}

// Customer chat interface
func PageChat(c *gin.Context) {
	referralSource := c.Query("refer") // More clear variable name

	if referralSource == "" {
		referralSource = c.Request.Referer()
	}
	if referralSource == "" {
		referralSource = "直接進入"
	}

	c.HTML(http.StatusOK, "chat_page.html", gin.H{
		"Refer": referralSource, // Keeping original template variable name
	})
}

// Support agent console
func PageChatMain(c *gin.Context) {
	c.HTML(http.StatusOK, "chat_main.html", nil)
}

// Mobile support agent console
func PageMobileChat(c *gin.Context) {
	c.HTML(http.StatusOK, "chat_mobile.html", nil)
}

// Settings
func PageSetting(c *gin.Context) {
	c.HTML(http.StatusOK, "setting.html", nil)
}
