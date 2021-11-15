package model

import "gorm.io/gorm"

type Secretshare struct {
	gorm.Model
	SecretId uint   `gorm:"type:int;not null"`
	UnitId   uint   `gorm:"type:int;not null"`
	Degree   int    `gorm:"type:int;not null"`
	Counter  int    `gorm:"type:int;not null"`
	Data     string `gorm:"type:string;not null"`
}
