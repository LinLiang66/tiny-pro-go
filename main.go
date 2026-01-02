package main

import (
	"fmt"
	"log"
	"tiny-admin-api-serve/middleware"
	jsonmiddleware "tiny-admin-api-serve/middleware/json"
	routers "tiny-admin-api-serve/routes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./config/config.yaml")

	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		panic(err)
		return
	}
	r := gin.Default()
	// 应用自定义JSON序列化中间件
	r.Use(jsonmiddleware.CustomJSON())
	// 应用全局鉴权中间件，默认所有路由都需要鉴权，只有标记了IsPublic的路由才开放
	r.Use(middleware.Auth.AuthRequired())
	// 注册路由
	routers.RouterUser(r)
	err = r.Run(":" + viper.GetString("port"))
	if err != nil {
		log.Printf("failed to start server: %v", err)
	}
}
