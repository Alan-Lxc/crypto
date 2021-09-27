package poly

import (
	"errors"
	"fmt"
	"github.com/ncw/gmp"
	"math/rand"
)

// Poly  p(x) = a0 + a1*x^1 + a2*x^2 + ... + an*x^n
//each polynomial is saved as p(x) = coef[0] + coef[1] *x^1+...+coef[n]*x^n
type Poly struct {
	coeff []*gmp.Int
}

func NewPoly(degree int) (Poly, error) {
	if degree < 0 {
		return Poly{}, errors.New(fmt.Sprintf("Can not get a poly with a negative degree"))
	}
	//A poly of n degree has  n+1 length
	coef := make([]*gmp.Int, degree+1)
	for i := 0; i < degree+1; i++ {
		coef[i] = gmp.NewInt(0)
	}
	return Poly{coef}, nil
}
func (poly Poly) GetDegree() int {
	deg := len(poly.coeff) - 1
	//Check the high end whether is 0,if so we should --
	for i := deg; i > 0; i-- {
		if poly.coeff[i].CmpInt32(0) == 0 {
			deg--
		} else {
			break
		}
	}
	return deg
}
func (poly Poly) GetCoeff(i int) (gmp.Int, error) {
	if i < 0 || i > len(poly.coeff)-1 {
		return *gmp.NewInt(0), errors.New("the parameter is out of range")
	}
	return *poly.coeff[i], nil
}
func (poly Poly) SetCoeffWithInt(i int, ci int64) error {
	if i < 0 || i > len(poly.coeff)-1 {
		return errors.New("the parameter is out of range")
	}
	poly.coeff[i].SetInt64(ci)
	return nil
}
func (poly Poly) SetCoeffWithGmp(i int, ci *gmp.Int) error {
	if i < 0 || i > len(poly.coeff)-1 {
		return errors.New("the parameter is out of range")
	}
	poly.coeff[i].Set(ci)
	return nil
}

//Reset the poly with coeff all equals 0
func (poly Poly) Reset() {
	for i := 0; i < len(poly.coeff); i++ {
		poly.coeff[i].SetInt64(0)
	}
}

//NewConstant create a poly p(x) = c
func NewConstant(c int64) Poly {
	poly, err := NewPoly(0)
	if err != nil {
		panic(err.Error())
	}
	poly.coeff[0] = gmp.NewInt(0)
	return poly
}

// NewRand returns a randomized polynomial with specified degree
// coefficients are pesudo-random numbers in [0, n)
func NewRand(degree int, rand *rand.Rand, n *gmp.Int) (Poly, error) {
	poly, err := NewPoly(degree)
	if err != nil {
		return Poly{}, err
	}
	poly.Rand(rand, n)
	return poly, nil
}
func (poly Poly) Rand(rand *rand.Rand, mod *gmp.Int) {

}
