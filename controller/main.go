package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"goflylivechat/common"
	"goflylivechat/models"
	"goflylivechat/tools"
	"goflylivechat/ws"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func PostInstall(c *gin.Context) {
	notExist, _ := tools.IsFileNotExist("./install.lock")
	if !notExist {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "系統已安裝",
		})
		return
	}
	server := c.PostForm("server")
	port := c.PostForm("port")
	database := c.PostForm("database")
	username := c.PostForm("username")
	password := c.PostForm("password")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, server, port, database)
	_, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
		tools.Logger().Println(err)
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "資料庫連線失敗：" + err.Error(),
		})
		return
	}
	isExist, _ := tools.IsFileExist(common.Dir)
	if !isExist {
		os.Mkdir(common.Dir, os.ModePerm)
	}
	fileConfig := common.MysqlConf
	file, _ := os.OpenFile(fileConfig, os.O_RDWR|os.O_CREATE, os.ModePerm)

	format := `{
	"Server":"%s",
	"Port":"%s",
	"Database":"%s",
	"Username":"%s",
	"Password":"%s"
}
`
	data := fmt.Sprintf(format, server, port, database, username, password)
	file.WriteString(data)
	models.Connect()
	installFile, _ := os.OpenFile("./install.lock", os.O_RDWR|os.O_CREATE, os.ModePerm)
	installFile.WriteString("gofly live chat")
	ok, err := install()
	if !ok {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "安裝成功",
	})
}
func install() (bool, error) {
	sqlFile := common.Dir + "go-fly.sql"
	isExit, _ := tools.IsFileExist(common.MysqlConf)
	dataExit, _ := tools.IsFileExist(sqlFile)
	if !isExit || !dataExit {
		return false, errors.New("config/mysql.json 資料庫設定檔或資料庫檔案 go-fly.sql 不存在")
	}
	sqls, _ := ioutil.ReadFile(sqlFile)
	sqlArr := strings.Split(string(sqls), "|")
	for _, sql := range sqlArr {
		if sql == "" {
			continue
		}
		err := models.Execute(sql)
		if err == nil {
			log.Println(sql, "\t success!")
		} else {
			log.Println(sql, err, "\t failed!")
		}
	}
	return true, nil
}

func GetStatistics(c *gin.Context) {
	visitors := models.CountVisitors()
	message := models.CountMessage(nil, nil)
	session := ws.VisitorCount()
	kefuNum := 0
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"visitors": visitors,
			"message":  message,
			"session":  session + kefuNum,
		},
	})
}
