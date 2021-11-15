package model

type Secret struct {
	secretname string `form:"secretname" json:"secretname" binding:"required"`
	degree     int    `form:"degree" json:"degree" binding:"required"`
	counter    int    `form:"counter" json:"counter" binding:"required"`
	secret     int    `form:"secret" json:"secret" binding:"required"`
}
