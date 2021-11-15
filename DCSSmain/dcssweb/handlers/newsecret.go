package handlers

import (
	"fmt"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

var secretid = 0

//dic = {}
type secret_data struct {
	Secretname string `form:"secretname" json:"secretname"`
	Secretid   string `form:"secretid" json:"secretid"`
	Degree     string `form:"degree" json:"degree"`
	Counter    string `form:"counter" json:"counter"`
	//Secret		string	`form:"secret" json:"secret"`
}

var Data = make([]secret_data, 0)
var controll *model.Controll

func Init_control() {
	controll = model.Initsystem()
}
func NewSecret(c *gin.Context) {
	secretname := c.PostForm("secretname")
	degree_s := c.PostForm("degree")
	counter_s := c.PostForm("counter")
	secret := c.PostForm("secret")
	// 查找该secretname是否存在。
	newsecret_check(secretname)
	//
	var tmp = secret_data{
		Secretname: secretname,
		Secretid:   strconv.Itoa(secretid),
		Degree:     degree_s,
		Counter:    counter_s,
		//Secret:   secret,
	}
	fmt.Println(tmp)
	Data = append(Data, tmp)
	secretid += 1
	//newsecret
	degree, _ := strconv.Atoi(degree_s)
	counter, _ := strconv.Atoi(counter_s)
	controll.NewSecret(secretid, degree, counter, secret)
}

func newsecret_check(secretname string) {
	//

}
