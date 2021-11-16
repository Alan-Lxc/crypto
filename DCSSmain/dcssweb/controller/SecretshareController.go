package controller

import (
	"github.com/Alan-Lxc/crypto_contest/dcssweb/common"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/model"
)

func Newsecretshare(secretId uint, unitId uint, degree int, counter int, data string) {
	db := common.GetDB()
	//	获取参数

	//向数据库中插入新纪录
	newSecretshareitem := model.Secretshare{
		SecretId: secretId,
		UnitId:   unitId,
		Degree:   degree,
		Counter:  counter,
		Data:     data,
	}
	db.Create(&newSecretshareitem)
	//返回结果

}
