package model

import (
	"github.com/Alan-Lxc/crypto_contest/dcssweb/systeminit"
	"strconv"
)

type Secret struct {
	secretname string `form:"secretname" json:"secretname" binding:"required""`
	degree     int    `form:"degree" json:"degree" binding:"required"`
	counter    int    `form:"counter" json:"counter" binding:"required"`
	secret     int    `form:"secret" json:"secret" binding:"required"`
}

var secretid = 0

func (s *Secret) Newsecret() {
	secretname := s.secretname
	degree := s.degree
	counter := s.counter
	secret := strconv.Itoa(s.secret)
	if secretname != "" {
		secretid += 1
	}
	systeminit.NewSecret(secretid, degree, counter, secret, systeminit.Controll{})
	return
}
