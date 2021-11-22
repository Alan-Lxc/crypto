package model

type SecretUnit struct {
	Secretid uint `gorm:"type:int;not null"`
	Unitid	uint `gorm:"type:int;not null"`
}