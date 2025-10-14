package main

import (
	"log"
	portal_page_restapi "portal_link/modules/portal_page/adapter/restapi"
	user_restapi "portal_link/modules/user/adapter/restapi"
	"portal_link/pkg"
	"portal_link/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置和資料庫連線
	config.Init()
	db := pkg.NewPG(config.GetDBConfig().DSN())
	defer db.Close()

	r := gin.Default()
	if err := user_restapi.NewUserHandler(r, db); err != nil {
		log.Fatal(err)
	}
	if err := portal_page_restapi.NewPortalPageHandler(r, db); err != nil {
		log.Fatal(err)
	}

	// 根路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Portal Link API Server is running!",
		})
	})

	// 啟動服務器
	r.Run(":8080")
}
