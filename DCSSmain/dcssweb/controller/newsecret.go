package controller

import (
	"github.com/Alan-Lxc/crypto_contest/dcssweb/common"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/model"
	reponse "github.com/Alan-Lxc/crypto_contest/dcssweb/response"
	"github.com/Alan-Lxc/crypto_contest/src/controller"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

//dic = {}
//type secret_data struct {
//	Secretname string `form:"secretname" json:"secretname"`
//	Secretid   string `form:"secretid" json:"secretid"`
//	Degree     string `form:"degree" json:"degree"`
//	Counter    string `form:"counter" json:"counter"`
//	//Secret		string	`form:"secret" json:"secret"`
//}

//var Data = make([]secret_data, 0)

//
//func Init_control() {
//	controll = model.Initsystem()
//}

func NewSecret(ctx *gin.Context) {
	db := common.GetDB()
	//	获取参数并进行数据验证
	secretname := ctx.PostForm("secretname")

	degree, err := strconv.Atoi(ctx.PostForm("degree"))
	if err != nil {
		reponse.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "t不符合规范")
		return
	}
	counter, err := strconv.Atoi(ctx.PostForm("counter"))
	if err != nil {
		reponse.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "n不符合规范")
		return
	}
	userId, err := strconv.Atoi(ctx.PostForm("userId"))
	if err != nil {
		reponse.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "userId不符合规范")
		return
	}
	Description := ctx.PostForm("description")
	secretcontent := ctx.PostForm("secret")
	//创建秘密
	newSecret := model.Secret{
		Secretname:    secretname,
		Degree:        int64(degree),
		Counter:       int64(counter),
		UserId:        uint(userId),
		Description:   Description,
		LastHandoffAt: time.Now(),
		Secret:        secretcontent,
	}
	db.Create(&newSecret)
	//newsecret
	controller.Controller.NewSecret(int(newSecret.ID), degree, counter, secretcontent)
}

func newsecret_check(secretname string) {
	//

}
