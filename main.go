package main

import (
	"fmt"
	"log"
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
	routers.RouterUser(r)
	err = r.Run(":" + viper.GetString("port"))
	if err != nil {
		log.Printf("failed to start server: %v", err)
	}
}
