package common

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func InitConfig() {
	workdDir, _ := os.Getwd()
	fmt.Println(workdDir)
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.SetConfigFile(workdDir + "/config/application.yml")
	err := viper.ReadInConfig()
	if err != nil {
		//fmt.Println(err.Error())
		panic(err)
	}
}
