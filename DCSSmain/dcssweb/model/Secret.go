package model

import (
	"gorm.io/gorm"
	"time"
)

//type Secret struct {
//	secretname string `form:"secretname" json:"secretname" binding:"required"`
//	degree     int    `form:"degree" json:"degree" binding:"required"`
//	counter    int    `form:"counter" json:"counter" binding:"required"`
//	secret     int    `form:"secret" json:"secret" binding:"required"`
//}
type Secret struct {
	gorm.Model
	Degree        int64  `gorm:"type:int;not null"`
	Counter       int64  `gorm:"type:int;not null"`
	UserId        uint   `gorm:"type:int;not null"`
	Description   string `gorm:"type:varchar(2000);not null"`
	Secretname    string `gorm:"type:varchar(200);"`
	Secret        string `gorm:"type:varchar(200)"`
	LastHandoffAt time.Time
}
