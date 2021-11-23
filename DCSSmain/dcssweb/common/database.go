package common

import (
	"fmt"
	model "github.com/Alan-Lxc/crypto_contest/dcssweb/model"
	srcmodel "github.com/Alan-Lxc/crypto_contest/src/model"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	//db, err := sql.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	if err != nil {
		panic("failed to connect to database, err: " + err.Error())
	}
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情

	err = db.AutoMigrate(&model.User{})
	err = db.AutoMigrate(&model.Unit{})
	err = db.AutoMigrate(&srcmodel.Secretshare{})
	err = db.AutoMigrate(&model.Secret{})
	//err = db.AutoMigrate(&model.Secretshare{})
	//err = db.Where("1=1").Delete(&model.Unit{}).Error
	//err = db.Where("1=1").Delete(&srcmodel.Secretshare{}).Error
	//err = db.Where("1=1").Delete(&model.Secret{}).Error
	if err != nil {
		return nil
	}
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
