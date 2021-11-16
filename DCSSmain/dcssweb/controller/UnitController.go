package controller

import (
	"github.com/Alan-Lxc/crypto_contest/dcssweb/common"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/model"
)

func NewUnit() {
	db := common.GetDB()

	newUnit := model.Unit{
		Secretsharenum: 0,
	}

	db.Create(&newUnit)
}
