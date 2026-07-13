package controller

import (
	"github.com/gin-gonic/gin"
	"goflylivechat/models"
	"goflylivechat/tools"
	"time"
)

// @Summary User Authentication API
// @Description Validates user credentials and returns access token
// @Tags Authentication
// @Produce json
// @Accept multipart/form-data
// @Param username formData string true "Registered username"
// @Param password formData string true "Account password"
// @Param type formData string true "Auth type (e.g., 'admin' or 'user')"
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /check [post]
func LoginCheckPass(c *gin.Context) {
	password := c.PostForm("password")
	username := c.PostForm("username")
	info := models.FindUser(username)

	// Authentication failed case
	if info.Name == "" || info.Password != tools.Md5(password) {
		c.JSON(200, gin.H{
			"code":    401,
			"message": "帳號或密碼不正確", // User-friendly message
		})
		return
	}

	// Prepare user session data
	userinfo := map[string]interface{}{
		"kefu_name":   info.Name,
		"kefu_id":     info.ID,
		"create_time": time.Now().Unix(),
	}

	// Token generation
	token, err := tools.MakeToken(userinfo)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "暫時無法登入，請稍後再試",
		})
		return
	}

	// Successful response
	c.JSON(200, gin.H{
		"code":    200,
		"message": "登入成功",
		"result": gin.H{
			"token":      token,
			"created_at": userinfo["create_time"],
		},
	})
}
