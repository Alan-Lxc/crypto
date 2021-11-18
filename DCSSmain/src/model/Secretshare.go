package model

import "gorm.io/gorm"

type Secretshare struct {
	gorm.Model
	SecretId uint   `gorm:"type:int;not null"`
	UnitId   uint   `gorm:"type:int;not null"`
	Degree   int    `gorm:"type:int;not null"`
	Counter  int    `gorm:"type:int;not null"`
	Row 		int 	`gorm:"type:int;not null"`
	Data     []byte `gorm:"type:varbinary(3000);not null"`
}
