package interpolation

import (
	. "github.com/Alan-Lxc/crypto_contest/basic/poly"
	"github.com/ncw/gmp"
)

//Get a polynomial that satisfy all x and y
func Lagrange(degree int, x []*gmp.Int, y []*gmp.Int, mod *gmp.Int) (Poly, error) {
	return Poly{}, nil
}
