package controller

import (
	"github.com/Alan-Lxc/crypto_contest/dcssweb/common"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/dto"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/model"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/response"
	"github.com/Alan-Lxc/crypto_contest/src/controller"
	model1 "github.com/Alan-Lxc/crypto_contest/src/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func HandoffSecret(ctx *gin.Context) {
	db := common.GetDB()

	secretid, err := strconv.Atoi(ctx.Query("secretid"))
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "id错误")
		return
	}
	var secret model.Secret
	result := db.Where("id = ?", secretid).First(&secret)
	if result.Error != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "查询错误")
		return
	}
	controller.Controller.Handoff(int(secret.ID), int(secret.Degree), int(secret.Counter))
	secret.ID = secret.ID
	db.Save(secret)
	//response
	response.Success(ctx, gin.H{}, "handoff success")
}
func GetSecretList(ctx *gin.Context) {
	db := common.GetDB()

	userid, err := strconv.Atoi(ctx.Query("userid"))
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "查询错误")
		return
	}
	var secrets []model.Secret
	db.Where(&model.Secret{UserId: uint(userid)}).Find(&secrets)
	gh := gin.H{}
	for _, v := range secrets {
		gh[strconv.Itoa(int(v.ID))] = v
	}

	response.Success(
		ctx,
		gin.H{
			"total":      len(secrets),
			"secretlist": secrets,
		},
		"获取成功",
	)

}
func GetSecret(ctx *gin.Context) {
	db := common.GetDB()

	//	获取参数并进行数据验证
	secretid, err := strconv.Atoi(ctx.Query("secretid"))
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "id错误")
		return
	}
	var secret model.Secret
	result := db.Where("id = ?", secretid).First(&secret)
	if result.Error != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "查询错误")
		return
	}
	//这里写重构秘密，把重构出的秘密
	response.Success(ctx, gin.H{
		"secret": dto.ToGetSecretDto(secret),
	}, "获取成功")
}

func ReconstructSecret(ctx *gin.Context) {
	db := common.GetDB()

	//	获取参数并进行数据验证
	secretid, err := strconv.Atoi(ctx.Query("secretid"))
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "id错误")
		return
	}
	var secret model.Secret
	result := db.Where("id = ?", secretid).First(&secret)
	if result.Error != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "查询错误")
		return
	}
	secretValue := controller.Controller.Reconstruct(int(secret.ID), int(secret.Degree), int(secret.Counter))
	response.Success(ctx, gin.H{
		"secret": secretValue,
	}, "获取成功")
}

func DeleteSecret(ctx *gin.Context) {
	db := common.GetDB()

	var secret model.Secret

	//	获取参数并进行数据验证
	id, err := strconv.Atoi(ctx.Query("secretid"))
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "id错误")
		return
	}

	db.Delete(&secret, id)
	db.Where("secret_id=?", id).Delete(&model1.Secretshare{})
}

//修改委员会成员数
func UpdateSecretCounter(ctx *gin.Context) {
	db := common.GetDB()

	//	获取参数并进行数据验证
	n, err := strconv.Atoi(ctx.Query("newcounter"))
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "counter不符合规范")
		return
	}
	secretid, err := strconv.Atoi(ctx.Query("secretid"))
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "id错误")
		return
	}
	var secret model.Secret
	db.First(&secret, secretid)
	controller.Controller.ModifyCommittee(int(secret.ID), int(secret.Degree), int(secret.Counter), n)
	secret.Counter = int64(n)
	db.Save(secret)
	response.Success(ctx, gin.H{}, "修改成功")

}

//func NewSecret2(ctx *gin.Context) {
//	db := common.GetDB()
//
//	//	获取参数并进行数据验证
//	secretname := ctx.PostForm("secretname")
//	degree, err := strconv.Atoi(ctx.PostForm("degree"))
//	if err != nil {
//		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "t不符合规范")
//		return
//	}
//	counter, err := strconv.Atoi(ctx.PostForm("counter"))
//	if err != nil {
//		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "n不符合规范")
//		return
//	}
//	userId, err := strconv.Atoi(ctx.PostForm("userId"))
//	if err != nil {
//		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "userId不符合规范")
//		return
//	}
//	Description := ctx.PostForm("description")
//	secretcontent := ctx.PostForm("secret")
//	//创建秘密
//	newSecret := model.Secret{
//		Secretname:  secretname,
//		Degree:      int64(degree),
//		Counter:     int64(counter),
//		UserId:      uint(userId),
//		Description: Description,
//		Secret:      secretcontent,
//	}
//	db.Create(&newSecret)
//	//id := newSecret.ID
//	//这里新建了秘密节点unit，要修改核心代码里的打印节点日志地址，并把地址存到下面的loglocation。
//	var tmpunits []model.Unit
//	res := db.Where("1=1").Find(&tmpunits)
//	currentunits := int(res.RowsAffected)
//	gap := counter-int(res.RowsAffected)
//	if gap>0{
//		for i := 0; i < gap; i++ {
//			tmpnum := strconv.Itoa(i+1+currentunits)
//			ip := "127.0.0.1:100"
//			if i+1+currentunits<10 {
//				ip+="0"
//			}
//			ip+=tmpnum
//			newunit := model.Unit{
//				UnitId: i+1+currentunits,
//				UnitIp: ip,
//			}
//			db.Create(&newunit)
//
//		}
//	}
//	res = db.Find(&tmpunits)
//	for i := 0; i < counter; i++ {
//		newsu := model.Secretshare{
//			Secretid: newSecret.ID,
//			Unitid:  uint(i+1),
//		}
//		db.Create(&newsu)
//
//	}
//
//	response.Success(ctx, gin.H{}, "新建秘密成功")
//	//secretId := newSecret.ID
//
//	//在前n个unit上分配秘密份额
//	//1. 在secretshareitem表中添加元素
//	//var unit model.Unit
//	//result:= db.Limit(counter).Find(&unit)
//	//if result.RowsAffected!=int64(counter) {
//	//	newUnitNum := counter-int(result.RowsAffected)
//	//	for i := 0; i < newUnitNum; i++ {
//	//		NewUnit()
//	//	}
//	//}
//	//result = db.Limit(counter).Select("ID").First(&unit)
//	//result.Get("ID")
//	//var data string
//	//for i := 0; i < counter; i++ {
//	////	将获取的秘密份额数组（长度为n），分别加入Secretshareitems表
//	//	data =
//	//	//	插入
//	//	Newsecretshareitem(int(secretId), counter, counter, data)
//	//}
//
//}
