package poly

import (
	"errors"
	"fmt"
	"github.com/ncw/gmp"
	"math/rand"
)

// Poly  p(x) = a0 + a1*x^1 + a2*x^2 + ... + an*x^n
//each polynomial is saved as p(x) = Coeffs[0] + Coeffs[1] *x^1+...+Coeffs[n]*x^n
type Poly struct {
	Coeffs []*gmp.Int
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
func (poly Poly) SetbyCoeff(coeff []*gmp.Int) {
	degree := len(coeff) - 1
	for i := 0; i <= degree; i++ {
		poly.SetCoeffWithGmp(i, coeff[i])
	}
}
func (poly Poly) GetDegree() int {
	deg := len(poly.Coeffs) - 1

	// note: i == 0 is not tested, because even the constant term is zero, we consider it's degree 0
	for i := deg; i > 0; i-- {
		//fmt.Println(poly.Coeffs[i])
		if poly.Coeffs[i].CmpInt32(0) == 0 {
			deg--
		} else {
			break
		}
	}

	return deg
}

// GetCoefficient returns Coeffs[i]
func (poly Poly) GetCoeff(i int) (gmp.Int, error) {
	if i < 0 || i >= len(poly.Coeffs) {
		return *gmp.NewInt(0), errors.New("out of boundary")
	}

	return *poly.Coeffs[i], nil
}
func (poly Poly) GetCoeffConstant() *gmp.Int {
	return poly.Coeffs[0]
}

// GetAllCoeffcients returns a copy of
func (poly Poly) GetAllCoeffs() (all []*gmp.Int) {
	all = make([]*gmp.Int, poly.GetDegree()+1)

	for i := range all {
		all[i] = gmp.NewInt(0)
		all[i] = poly.Coeffs[i]
	}

	return all
}
func (poly Poly) GetPtrtoCoeff(i int) (*gmp.Int, error) {
	return poly.Coeffs[i], nil
}

func (poly Poly) GetPtrtoConstant() *gmp.Int {
	return poly.Coeffs[0]
}
func (poly *Poly) SetCoeffWithInt(i int, ci int64) error {
	if i < 0 || i > len(poly.Coeffs)-1 {
		return errors.New("the parameter is out of range")
	}
	poly.Coeffs[i].SetInt64(ci)
	return nil
}
func (poly *Poly) SetCoeffWithGmp(i int, ci *gmp.Int) error {
	if i < 0 || i > len(poly.Coeffs)-1 {
		return errors.New("the parameter is out of range")
	}
	poly.Coeffs[i].Set(ci)
	return nil
}

//Reset the poly with Coeffs all equals 0
func (poly *Poly) Reset() {
	for i := 0; i < len(poly.Coeffs); i++ {
		poly.Coeffs[i].SetInt64(0)
	}
}

//NewConstant create a poly p(x) = c
func NewConstant(c int64) Poly {
	poly, err := NewPoly(0)
	if err != nil {
		panic(err.Error())
	}
	poly.Coeffs[0] = gmp.NewInt(c)
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
	for i := range poly.Coeffs {
		poly.Coeffs[i].Rand(rand, mod)
	}

	highest := len(poly.Coeffs) - 1

	for {
		if 0 == poly.Coeffs[highest].CmpInt32(0) {
			poly.Coeffs[highest].Rand(rand, mod)
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
	if degree+1 <= len(poly.Coeffs) {
		poly.Coeffs = poly.Coeffs[:degree+1]
	} else {
		//or we need to grow the size
		extra := make([]*gmp.Int, degree+1-len(poly.Coeffs))
		for i := 0; i < len(extra); i++ {
			extra[i] = gmp.NewInt(0)
		}
		poly.Coeffs = append(poly.Coeffs, extra...)

	}
	poly.Reset()
	return nil
}
func (poly *Poly) ResetTo(other Poly) {
	err := poly.ResetDegree(other.GetDegree())
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < other.GetDegree()+1; i++ {
		poly.Coeffs[i].Set(other.Coeffs[i])
	}
}

func (poly *Poly) Add(poly1 Poly, poly2 Poly) {
	deg1 := poly1.GetDegree()
	deg2 := poly2.GetDegree()

	if deg1 > deg2 {
		poly.ResetTo(poly1)
	} else {
		poly.ResetTo(poly2)
	}

	for i := 0; i < min(deg1, deg2)+1; i++ {
		poly.Coeffs[i].Add(poly1.Coeffs[i], poly2.Coeffs[i])
	}

	poly.Coeffs = poly.Coeffs[:poly.GetDegree()+1]

	//let the poly as long as the longest (highest end =longest)
	//and then add poly1 + poly2
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func (poly *Poly) Sub(poly1 Poly, poly2 Poly) {
	// make sure poly is as long as the longest of op1 and op2
	deg1 := poly1.GetDegree()
	deg2 := poly2.GetDegree()

	if deg1 > deg2 {
		poly.ResetTo(poly1)
	} else {
		poly.ResetTo(poly2)
	}
	//fmt.Println("deg-1")
	for i := 0; i < min(deg1, deg2); i++ {
		poly.Coeffs[i].Sub(poly1.Coeffs[i], poly2.Coeffs[i])
	}

	poly.Coeffs = poly.Coeffs[:poly.GetDegree()+1]

	//if poly1.GetDegree() > poly2.GetDegree() {
	//	poly.ResetTo(poly1)
	//	for i := 0; i < poly2.GetDegree(); i++ {
	//		poly.Coeffs[i].Sub(poly1.Coeffs[i], poly2.Coeffs[i])
	//	}
	//} else {
	//	poly.ResetTo(poly2)
	//	for i := 0; i < poly1.GetDegree(); i++ {
	//		poly.Coeffs[i].Sub(poly1.Coeffs[i], poly2.Coeffs[i])
	//	}
	//
	//} //let the poly as long as the longest (highest end =longest)
	////and then add poly1 - poly2

}
func (poly *Poly) Multiply(poly1 Poly, poly2 Poly) error {
	deg1 := poly1.GetDegree()
	deg2 := poly2.GetDegree()
	err := poly.ResetDegree(deg1 + deg2)
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i <= deg1; i++ {
		for j := 0; j <= deg2; j++ {
			poly.Coeffs[i+j].AddMul(poly1.Coeffs[i], poly2.Coeffs[j])
		}
	}
	poly.Coeffs = poly.Coeffs[:poly.GetDegree()+1]
	return err
}

func (poly Poly) DeepCopy() Poly {
	dst, err := NewPoly(poly.GetDegree())
	if err != nil {
		panic("deepcopy failed: " + err.Error())
	}

	for i := 0; i < len(dst.Coeffs); i++ {
		dst.Coeffs[i].Set(poly.Coeffs[i])
	}

	return dst
}

// MulSelf set poly to poly * op
func (poly *Poly) MulSelf(op Poly) error {
	op1 := poly.DeepCopy()
	poly.Multiply(op1, op)
	return nil
}

//AddMul  return to self + k * poly2
func (poly *Poly) AddMul(poly2 Poly, k *gmp.Int) {
	for i := 0; i <= poly2.GetDegree(); i++ {
		poly.Coeffs[i].AddMul(poly2.Coeffs[i], k)
	}
}

func (poly *Poly) Mod(p *gmp.Int) {
	for i := 0; i < len(poly.Coeffs); i++ {
		poly.Coeffs[i].Mod(poly.Coeffs[i], p)
	}
}

//求解多项式的值
// EvalMod returns poly(x) using Horner's rule. If p != nil, returns poly(x) mod p
func (poly Poly) EvalMod(x *gmp.Int, p *gmp.Int, result *gmp.Int) {
	result.Set(poly.Coeffs[poly.GetDegree()])

	for i := poly.GetDegree(); i >= 1; i-- {
		result.Mul(result, x)
		result.Add(result, poly.Coeffs[i-1])
	}

	if p != nil {
		result.Mod(result, p)
	}
}

//op2一定是一次多项式形式
// Divide sets poly to op1 / op2. **op2 must be of format x+a **
//
// Complexity is O(deg1)
func (poly *Poly) Divide(op1 Poly, op2 Poly) error {
	degree1 := op1.GetDegree()
	degree2 := op2.GetDegree()

	poly.ResetDegree(degree1)

	if degree2 != 1 {
		return errors.New("op2 must be of format x-a")
	}

	if len(poly.Coeffs) < degree1-1 {
		return errors.New("receiver too small")
	}

	tmp := gmp.NewInt(0)

	numerator, err := NewPoly(degree1)
	if err != nil {
		return errors.New("unknown error")
	}

	for i := 0; i <= degree1; i++ {
		numerator.Coeffs[i].Set(op1.Coeffs[i])
	}

	for i := degree1; i > 0; i-- {
		poly.Coeffs[i-1].Div(numerator.Coeffs[i], op2.Coeffs[degree2])
		for j := degree2; j >= 0; j-- {
			tmp.Mul(poly.Coeffs[i-1], op2.Coeffs[j])
			numerator.Coeffs[i+j-degree2].Sub(numerator.Coeffs[i+j-degree2], tmp)
		}
	}
	p := gmp.NewInt(0)
	p.SetString("57896044618658097711785492504343953926634992332820282019728792006155588075521", 10)
	//fmt.Println("coeffs0 is ", gmp.NewInt(0).Mod(numerator.Coeffs[0], p))
	poly.Coeffs = poly.Coeffs[:poly.GetDegree()+1]
	return nil
}

func (poly Poly) Copy() Poly {
	tmp, _ := NewPoly(poly.GetDegree())

	for i := 0; i < len(tmp.Coeffs); i++ {
		tmp.Coeffs[i].Set(poly.Coeffs[i])
	}

	return tmp
}
func (poly Poly) GetAllCoeff() []*gmp.Int {
	res := make([]*gmp.Int, poly.GetDegree()+1)
	for i := range res {
		res[i] = gmp.NewInt(0)
		res[i] = poly.Coeffs[i]
	}
	return res
}
func FromVec(coeff ...int64) Poly {
	if len(coeff) == 0 {
		return NewConstant(0)
	}
	deg := len(coeff) - 1

	poly, err := NewPoly(deg)
	if err != nil {
		panic(err.Error())
	}

	for i := range poly.Coeffs {
		poly.Coeffs[i].SetInt64(coeff[i])
	}

	return poly
}
func NewEmpty() Poly {
	return NewConstant(0)
}
func (poly Poly) IsSame(op Poly) bool {
	if op.GetDegree() != poly.GetDegree() {
		return false
	}

	for i := 0; i <= op.GetDegree(); i++ {
		if op.Coeffs[i].Cmp(poly.Coeffs[i]) != 0 {
			return false
		}
	}

	return true
}

func (poly Poly) GetCap() int {
	return len(poly.Coeffs)
}
func (poly *Poly) GrowCapTo(cap int) {
	current := poly.GetCap()
	if cap <= current {
		return
	}

	// if we need to grow the slice
	needed := cap - current
	neededPointers := make([]*gmp.Int, needed)
	for i := 0; i < len(neededPointers); i++ {
		neededPointers[i] = gmp.NewInt(0)
	}

	poly.Coeffs = append(poly.Coeffs, neededPointers...)
}

// AddSelf sets poly to poly + op
func (poly *Poly) AddSelf(op Poly) {
	op1 := poly.DeepCopy()
	poly.Add(op1, op)
}

// SubSelf sets poly to poly - op
func (poly *Poly) SubSelf(op Poly) error {
	// make sure poly is as long as the longest of op1 and op2
	deg1 := op.GetDegree()

	poly.GrowCapTo(deg1 + 1)

	for i := 0; i < deg1+1; i++ {
		fmt.Println(poly.Coeffs[i], op.Coeffs[i])
		poly.Coeffs[i].Sub(poly.Coeffs[i], op.Coeffs[i])
	}

	poly.Coeffs = poly.Coeffs[:poly.GetDegree()+1]

	// FIXME: no need to return error
	return nil
}

// IsZero returns if poly == 0
func (poly Poly) IsZero() bool {
	if poly.GetDegree() != 0 {
		return false
	}

	return poly.GetPtrtoConstant().CmpInt32(0) == 0
}
func (poly Poly) GetLeadingCoefficient() gmp.Int {
	lc := gmp.NewInt(0)
	lc.Set(poly.Coeffs[poly.GetDegree()])

	return *lc
}

// DivMod sets computes q, r such that a = b*q + r.
// This is an implementation of Euclidean division. The complexity is O(n^3)!!
func DivMod(a Poly, b Poly, p *gmp.Int, q, r *Poly) (err error) {
	if b.IsZero() {
		return errors.New("divide by zero")
	}

	q.ResetDegree(0)
	r.ResetTo(a)

	d := b.GetDegree()
	c := b.GetLeadingCoefficient()

	// cInv = 1/c
	cInv := gmp.NewInt(0)
	cInv.ModInverse(&c, p)

	for r.GetDegree() >= d {
		lc := r.GetLeadingCoefficient()
		s, err := NewPoly(r.GetDegree() - d)
		if err != nil {
			return err
		}

		s.SetCoeffWithGmp(r.GetDegree()-d, lc.Mul(&lc, cInv))

		q.AddSelf(s)

		sb := NewEmpty()
		sb.Multiply(s, b)

		// deg r reduces by each iteration
		r.SubSelf(sb)

		// modulo p
		q.Mod(p)
		r.Mod(p)
	}

	return nil
}
