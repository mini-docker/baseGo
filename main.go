package main

import (
	"baseGo/common"
	"baseGo/routes"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func main() {
	InitConfig()
	db := common.InitDB()
	defer db.Close()

	r := gin.Default()
	r = routes.CollectRoute(r)
	port := viper.GetString("server.port")
	// if port != "" {
	// 	panic(r.Run(":", port))
	// }
	// panic(r.Run())
	r.Run()
	fmt.Println(port)
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
