package interpolation

import (
	. "github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	"github.com/ncw/gmp"
)

//Get a polynomial that satisfy all x and y
//degree
func Lagrange(degree int, x []*gmp.Int, y []*gmp.Int, mod *gmp.Int) (Poly, error) {
	return Poly{}, nil

}
