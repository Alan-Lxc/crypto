package handlers

import (
	"fmt"
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

//func Getstart() {
//	file = "./data"
//
//
//}
func NewSecret(c *gin.Context) {
	secretname := c.PostForm("secretname")
	degree := c.PostForm("degree")
	counter := c.PostForm("counter")
	//secret := c.PostForm("secret")
	// 查找该secretname是否存在。
	newsecret_check(secretname)

	//
	var tmp = secret_data{
		Secretname: secretname,
		Secretid:   strconv.Itoa(secretid),
		Degree:     degree,
		Counter:    counter,
		//Secret:   secret,
	}
	fmt.Println(tmp)
	//jsonbyte,err := json.Marshal(tmp)
	//
	Data = append(Data, tmp)
	secretid += 1
	//newsecret

}

func newsecret_check(secretname string) {
	//

}
