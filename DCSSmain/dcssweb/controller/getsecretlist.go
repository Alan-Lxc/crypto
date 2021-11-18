package controller
////
//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"net/http"
//)
//
//type secretlist struct {
//	total int
//	data  []secret_data
//}
//
//func Getsecretlist(c *gin.Context) {
//
//	//get secret list from database
//	if len(Data) != 0 {
//		fmt.Println(gin.H{
//			"total":    len(Data),
//			"datalist": Data,
//		})
//		c.JSON(http.StatusOK, gin.H{
//			"total":    len(Data),
//			"datalist": Data,
//		})
//		//data1,_ := json.Marshal(Data)
//		//fmt.Println(data1)
//	}
//
//}
