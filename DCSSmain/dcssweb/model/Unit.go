package model

import "gorm.io/gorm"

type Unit struct {
	gorm.Model
	UnitId int `gorm:"type:int;not null"`
	UnitIp string `gorm:"type:varchar(100);not null"`
	//Secretnum int64  `gorm:"type:int;not null"`
	//Loglocation    string `gorm:"type:varchar(200);not null"`
}

