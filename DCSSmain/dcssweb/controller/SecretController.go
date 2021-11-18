package controller

import (
	"fmt"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/common"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/dto"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/model"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var secret model.Secret
var secrets []model.Secret

func RetrieveSecretByUserid(ctx *gin.Context) {
	db := common.GetDB()

	userid, err := strconv.Atoi(ctx.Query("userid"))
	if err != nil {
		reponse.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "查询错误")
		return
	}
	db.Where(&model.Secret{UserId: uint(userid)}).Find(&secrets)
	fmt.Println(secrets)
	gh := gin.H{}
	for _, v := range secrets {
		gh[strconv.Itoa(int(v.ID))] = v
	}
	reponse.Success(ctx, gh, "获取成功")

}
func RetrieveSecretById(ctx *gin.Context) {
	db := common.GetDB()

	//	获取参数并进行数据验证
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		reponse.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "id错误")
		return
	}

	result := db.Where("user_id = ?", id).Find(&secret)
	if result.Error != nil {
		reponse.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "查询错误")
		return
	}
	//这里写重构秘密，把重构出的秘密
	reponse.Success(ctx, gin.H{
		"secret": dto.ToRetrieveSecretByIdDto(secret),
	}, "获取成功")
}
func ReconstructSecret(ctx *gin.Context) {
	RetrieveSecretById(ctx)
}
func DeleteSecret(ctx *gin.Context) {
	db := common.GetDB()

	var secret model.Secret
	//	获取参数并进行数据验证
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		reponse.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "id错误")
		return
	}

	db.Delete(&secret, id)
}

//修改委员会成员数
func UpdateSecretCounter(ctx *gin.Context) {
	db := common.GetDB()

	//	获取参数并进行数据验证
	n, err := strconv.Atoi(ctx.PostForm("counter"))
	if err != nil {
		reponse.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "counter不符合规范")
		return
	}
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		reponse.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "id错误")
		return
	}

	db.First(&secret, id)
	secret.Counter = int64(n)
	db.Save(secret)

	reponse.Success(ctx, gin.H{}, "修改成功")

}

//func NewSecret(ctx *gin.Context) {
//	db := common.GetDB()
//
//	//	获取参数并进行数据验证
//	secretname := ctx.PostForm("secretname")
//
//	degree, err := strconv.Atoi(ctx.PostForm("degree"))
//	if err != nil {
//		reponse.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "t不符合规范")
//		return
//	}
//	counter, err := strconv.Atoi(ctx.PostForm("counter"))
//	if err != nil {
//		reponse.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "n不符合规范")
//		return
//	}
//	userId, err := strconv.Atoi(ctx.PostForm("userId"))
//	if err != nil {
//		reponse.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "userId不符合规范")
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
//	id := newSecret.ID
//	//这里新建了秘密节点unit，要修改核心代码里的打印节点日志地址，并把地址存到下面的loglocation。
//	for i := 0; i < counter; i++ {
//		newunit := model.Unit{
//			Secretsharenum: 1,
//			Loglocation:    "",
//		}
//		Newsecretshare(id, newunit.ID, degree, counter, "")
//		db.Create(&newunit)
//	}
//	reponse.Success(ctx, gin.H{}, "新建秘密成功")
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
