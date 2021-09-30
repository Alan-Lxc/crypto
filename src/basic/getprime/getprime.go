package getprime

import (
	"crypto/rand"
	"github.com/ncw/gmp"
)

func Gcd(a, b *gmp.Int) *gmp.Int {
	if b.CmpInt32(0) == 1 {
		return a
	} else {
		tmp := gmp.NewInt(0)
		tmp = tmp.Mod(a, b)
		return Gcd(b, tmp)
	}
}

func MillerRabin(x *gmp.Int) bool {
	if x.CmpInt32(2) == -1 {
		return false
	}
	if x.CmpInt32(2) == 0 {
		return true
	}
	q := gmp.NewInt(0)
	q = q.Sub(x, gmp.NewInt(1))
	k := 0
	tmp := gmp.NewInt(2)
	for true {
		tmp = tmp.Mod(q, gmp.NewInt(2))
		if tmp.CmpInt32(0) == 0 && q.CmpInt32(0) == 1 {
			break
		}
		q = q.Div(q, gmp.NewInt(2))
		k += 1
	}
	for count := 100; count > 0; count-- {
		//t := rand.Int63n(int64(tmp.Sub(x,gmp.NewInt(1))))
		//if Gcd(gmp.NewInt(int64(t),)
	}
	return true
}
func GetPrime(bits int) *gmp.Int {
	r, _ := rand.Prime(rand.Reader, bits)
	tmp := gmp.NewInt(0)
	tmp.SetString(r.String(), 10)
	return tmp
}
