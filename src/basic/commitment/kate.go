package commitment

import (
	"fmt"
	"github.com/Alan-Lxc/crypto_contest/src/basic/conv"
	"github.com/Alan-Lxc/crypto_contest/src/basic/ecparam"
	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	. "github.com/Nik-U/pbc"
	. "github.com/ncw/gmp"
	"math/big"
)

type DLPolyCommit struct {
	pairing *Pairing
	pk      []*Power
	degree  int
	p       *Int
}

// Generate New G1
func (c *DLPolyCommit) NewG1() *Element {
	return c.pairing.NewG1()
}

//Generate New GT
func (c *DLPolyCommit) NewGT() *Element {
	return c.pairing.NewGT()
}

// polyEval sets res to poly(x)
func (c *DLPolyCommit) polyEval(res *Int, poly poly.Poly, x *Int) {

	poly.EvalMod(x, c.p, res)
}

// Let poly(x)=c0 + c1*x + ... cn * x^n, polyEvalInExponent sets res to g^poly(alpha)
func (c *DLPolyCommit) polyEvalInExponent(res *Element, poly poly.Poly) {
	// res = 1
	res.Set1()
	tmp := c.pairing.NewG1()
	//fmt.Println(poly)
	for i := 0; i <= poly.GetDegree(); i++ {
		//fmt.Println(i  ,111)
		// tmp = g^{a^i} ^ ci
		ci, err := poly.GetCoeff(i)
		if err != nil {
			panic("can't get coeff i")
		}

		//fmt.Println(i  ,222)
		//fmt.Printf("%d\n",len(c.pk))
		//fmt.Printf("%d\n",poly.GetDegree())
		//fmt.Println(i,"gottt",conv.GmpInt2BigInt(&ci))
		//c.pk[i].PowBig(tmp, conv.GmpInt2BigInt(&ci))
		//b := new(big.Int)
		//b.SetBytes((&ci).Bytes())
		//fmt.Println(i  ,333)
		//b = big.NewInt(32)
		//tmp.PowBig(c.pk[i],b)
		//fmt.Println((c.pk[i]).Source(),b,c.pk[i].PowBig(tmp, b))
		//c.pk[i].PowBig(tmp, b)
		c.pk[i].PowBig(tmp, conv.GmpInt2BigInt(&ci))
		//fmt.Println(i  ,444)
		//fmt.Println(i,"gottt",conv.GmpInt2BigInt(&ci))
		res.Mul(res, tmp)
	}
}

// print the public keys
func (c *DLPolyCommit) printPublicKey() {
	for i := 0; i <= c.degree; i++ {
		fmt.Printf("g^(SK^%d): %s\n", i, c.pk[i].Source().String())
	}
}

var Curve = ecparam.PBC256

// SetupFix initializes a fixed pairing
func (c *DLPolyCommit) SetupFix(degree int) {
	c.degree = degree

	// setup the pairing
	c.pairing = Curve.Pairing
	c.p = Curve.Gmp

	// trusted setup
	c.pk = make([]*Power, degree+1)

	// a generator g
	g := Curve.Element

	// secret key
	sk := new(big.Int)
	sk.SetString("1", 10)

	tmp := new(big.Int)
	for i := 0; i <= degree; i++ {
		// tmp = sk ^ i
		bigP := big.NewInt(0)
		bigP.SetString(c.p.String(), 10)
		tmp.Exp(sk, big.NewInt(int64(i)), bigP)
		//tmp.Exp(sk, big.NewInt(int64(i)),nil)
		// pk[i] = g ^ tmp
		// Search pk and modify them all
		inter := c.pairing.NewG1()
		c.pk[i] = inter.PowBig(g, tmp).PreparePower()
	}
}

// Commit sets res to g^poly(alpha)
func (c *DLPolyCommit) Commit(res *Element, poly poly.Poly) {
	c.polyEvalInExponent(res, poly)
}

// Open is not used
func (c *DLPolyCommit) Open() {
	panic("unimplemented")
}

// VerifyPoly checks C == g ^ poly(alpha)
func (c *DLPolyCommit) VerifyPoly(C *Element, poly poly.Poly) bool {
	tmp := c.pairing.NewG1()
	c.polyEvalInExponent(tmp, poly)
	return tmp.Equals(C)
}

// CreateWitness sets res to g ^ phi(alpha) where phi(x) = (poly(x)-poly(x0)) / (x - x0)
func (c *DLPolyCommit) CreateWitness(res *Element, polynomial poly.Poly, x0 *Int) {
	poly_t := polynomial.DeepCopy()

	// tmp = polynomial(x0)
	tmp := new(Int)
	c.polyEval(tmp, poly_t, x0)
	// fmt.Printf("CreateWitness\n%s\n%s\n", polynomial.String(), tmp.String())

	// poly_t = polynomial(x)-polynomial(x0)
	poly_t.GetPtrtoConstant().Sub(poly_t.GetPtrtoConstant(), tmp)
	//fmt.Println("did sub")
	// quot == poly_t / (x - x0)
	quot := poly.NewEmpty()

	// denominator = x - x0
	denominator, err := poly.NewPoly(1)
	if err != nil {
		panic("can't create poly")
	}
	// FIXME: converting to int64 is dangerous
	denominator.SetCoeffWithInt(1, 1)
	denominator.GetPtrtoConstant().Neg(x0)

	quot.Divide(poly_t, denominator)
	// fmt.Printf("CreateWitness2\n%s\n", quot.String())
	//fmt.Println("t1",time.Now())
	c.polyEvalInExponent(res, quot)
	//fmt.Println("t2",time.Now())
}

// VerifyEval checks the correctness of w, returns true/false
func (c *DLPolyCommit) VerifyEval(C *Element, x *Int, polyX *Int, w *Element) bool {
	//test
	//d1 := c.pairing.NewGT()
	//d2 := c.pairing.NewGT()
	//d1.Pair(c.pk[0].Source(), c.pk[0].Source())
	//d2.Pair(c.pk[0].Source(), c.pk[0].Source())
	//dd := big.NewInt(0)
	//dd.SetString(c.p.String(), 10)
	//d1.PowBig(d1, dd)
	//d2.PowBig(d2, big.NewInt(0).Mul(dd, dd))
	//fmt.Println("?????????????????", d1.Equals(d2))

	e1 := c.pairing.NewGT()
	e2 := c.pairing.NewGT()
	t1 := c.pairing.NewGT()
	t2 := c.pairing.NewG1()
	e1.Pair(C, c.pk[0].Source())
	exp := big.NewInt(0)
	exp.SetString(x.String(), 10)
	c.pk[0].PowBig(t2, exp)
	t2.Div(c.pk[1].Source(), t2)
	e2.Pair(w, t2)
	t1.Pair(c.pk[0].Source(), c.pk[0].Source())
	exp.SetString(polyX.String(), 10)
	t1.PowBig(t1, exp)
	e2.Mul(e2, t1)
	// fmt.Printf("e1\n%s\ne2\n%s\n", e1.String(), e2.String())
	return e1.Equals(e2)
}

func (c *DLPolyCommit) CalcAmtWitness(C *Element, Witness, tmpWitness []*Element, polyX *Int, step int) bool {

	e1 := c.pairing.NewGT()
	e2 := c.pairing.NewGT()
	e1.Pair(c.pk[0].Source(), c.pk[0].Source())
	for i := 0; i < step; i++ {
		e2.Pair(Witness[i], tmpWitness[i])
		e1.Mul(e1, e2)
	}
	e2.Pair(c.pk[0].Source(), c.pk[0].Source())
	exp := big.NewInt(0)
	exp.SetString(polyX.String(), 10)
	e2.PowBig(e2, exp)
	e1.Mul(e1, e2)
	e2.Pair(C, c.pk[0].Source())
	return e2.Equals(e1)
}
