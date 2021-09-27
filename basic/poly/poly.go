package poly

import (
	"errors"
	"fmt"
	"github.com/ncw/gmp"
	"math/rand"
)

// Poly  p(x) = a0 + a1*x^1 + a2*x^2 + ... + an*x^n
//each polynomial is saved as p(x) = coeff[0] + coeff[1] *x^1+...+coeff[n]*x^n
type Poly struct {
	coeff []*gmp.Int
}

func NewPoly(degree int) (Poly, error) {
	if degree < 0 {
		return Poly{}, errors.New(fmt.Sprintf("Can not get a poly with a negative degree"))
	}
	//A poly of n degree has  n+1 length
	coeff := make([]*gmp.Int, degree+1)
	for i := 0; i < degree+1; i++ {
		coeff[i] = gmp.NewInt(0)
	}
	return Poly{coeff}, nil
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
func (poly *Poly) SetCoeffWithInt(i int, ci int64) error {
	if i < 0 || i > len(poly.coeff)-1 {
		return errors.New("the parameter is out of range")
	}
	poly.coeff[i].SetInt64(ci)
	return nil
}
func (poly *Poly) SetCoeffWithGmp(i int, ci *gmp.Int) error {
	if i < 0 || i > len(poly.coeff)-1 {
		return errors.New("the parameter is out of range")
	}
	poly.coeff[i].Set(ci)
	return nil
}

//Reset the poly with coeff all equals 0
func (poly *Poly) Reset() {
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
	poly.coeff[0] = gmp.NewInt(c)
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
	for i := range poly.coeff {
		poly.coeff[i].Rand(rand, mod)
	}

	highest := len(poly.coeff) - 1

	for {
		if 0 == poly.coeff[highest].CmpInt32(0) {
			poly.coeff[highest].Rand(rand, mod)
		} else {
			break
		}
	}
}
func (poly *Poly) ResetDegree(degree int) error {
	if degree < 0 {
		return errors.New("the parameter is out of range")
	}
	//want degree less than already,need to shrink the degree
	if degree+1 <= len(poly.coeff) {
		poly.coeff = poly.coeff[:degree+1]
	} else {
		//or we need to grow the size
		extra := make([]*gmp.Int, degree+1-len(poly.coeff))
		for i := 0; i < len(extra); i++ {
			extra[i] = gmp.NewInt(0)
		}
		poly.coeff = append(poly.coeff, extra...)
	}
	return nil
}
func (poly *Poly) ResetTo(other Poly) {
	err := poly.ResetDegree(other.GetDegree())
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < other.GetDegree()+1; i++ {
		poly.coeff[i].Set(other.coeff[i])
	}
}

func (poly Poly) Add(poly1 Poly, poly2 Poly) {
	if poly1.GetDegree() > poly2.GetDegree() {
		poly.ResetTo(poly1)
		for i := 0; i < poly2.GetDegree(); i++ {
			poly.coeff[i].Add(poly1.coeff[i], poly2.coeff[i])
		}
	} else {
		poly.ResetTo(poly2)
		for i := 0; i < poly1.GetDegree(); i++ {
			poly.coeff[i].Add(poly1.coeff[i], poly2.coeff[i])
		}

	} //let the poly as long as the longest (highest end =longest)
	//and then add poly1 + poly2
}
func (poly Poly) Sub(poly1 Poly, poly2 Poly) {
	if poly1.GetDegree() > poly2.GetDegree() {
		poly.ResetTo(poly1)
		for i := 0; i < poly2.GetDegree(); i++ {
			poly.coeff[i].Sub(poly1.coeff[i], poly2.coeff[i])
		}
	} else {
		poly.ResetTo(poly2)
		for i := 0; i < poly1.GetDegree(); i++ {
			poly.coeff[i].Sub(poly1.coeff[i], poly2.coeff[i])
		}

	} //let the poly as long as the longest (highest end =longest)
	//and then add poly1 - poly2

}
func (poly Poly) Multiply(poly1 Poly, poly2 Poly) {
	deg1 := poly1.GetDegree()
	deg2 := poly2.GetDegree()
	err := poly.ResetDegree(deg1 + deg2)
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i <= deg1; i++ {
		for j := 0; i <= deg2; j++ {
			poly.coeff[i+j].AddMul(poly1.coeff[i], poly2.coeff[j])
		}
	}
	poly.coeff = poly.coeff[:poly.GetDegree()+1]

}

//AddMul  return to self + k * poly2
func (poly Poly) AddMul(poly2 Poly, k *gmp.Int) {
	for i := 0; i <= poly2.GetDegree(); i++ {
		poly.coeff[i].AddMul(poly2.coeff[i], k)
	}
}

func (poly Poly) Mod(p *gmp.Int) {
	for i := 0; i < len(poly.coeff); i++ {
		poly.coeff[i].Mod(poly.coeff[i], p)
	}
}
