package controller

import (
	"github.com/gin-gonic/gin"
	"goflylivechat/models"
)

func GetRoleList(c *gin.Context) {
	roles := models.FindRoles()
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "取得成功",
		"result": roles,
	})
}
func PostRole(c *gin.Context) {
	roleId := c.PostForm("id")
	method := c.PostForm("method")
	name := c.PostForm("name")
	path := c.PostForm("path")
	if roleId == "" || method == "" || name == "" || path == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "參數不能為空",
		})
		return
	}
	models.SaveRole(roleId, name, method, path)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}
