package main

import (
	"baseGo/common"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func main() {
	InitConfig()
	db := common.InitDB()
	defer db.Close()

	port := viper.GetString("server.port")
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
