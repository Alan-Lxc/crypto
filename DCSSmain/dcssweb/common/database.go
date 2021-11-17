package common

import (
	"fmt"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/model"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	host := "127.0.0.1"
	port := "3306"
	database := "backend"
	username := "root"
	password := "root"
	charset := "utf8"
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
	err = db.AutoMigrate(&model.Secretshare{})
	err = db.AutoMigrate(&model.Secret{})
	if err != nil {
		return nil
	}
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
