package polycommitment

import (
	"bytes"
	"encoding/gob"
	"github.com/Alan-Lxc/crypto_contest/src/basic/conv"
	"github.com/Alan-Lxc/crypto_contest/src/basic/ecparam"
	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	"github.com/Nik-U/pbc"
	"math/big"
)

var Curve = ecparam.PBC256

type PolyCommit struct {
	c []*pbc.Element
}

func (com PolyCommit) Equals(b PolyCommit) bool {
	if len(b.c) != len(com.c) {
		return false
	}
	for i := range b.c {
		if !com.c[i].Equals(b.c[i]) {
			return false
		}
	}
	return true
}
func (com PolyCommit) Bytes() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	binary := make([][]byte, len(com.c))

	for i := range binary {
		binary[i] = com.c[i].CompressedBytes()
		if com.c[i].Is0() {
			binary[i][0] = 0xff
		}
	}

	if err := enc.Encode(binary); err != nil {
		return nil
	}
	return buf.Bytes()
}
func NewPolyCommit(poly poly.Poly) PolyCommit {
	allCoef := poly.GetAllCoeff()
	C := make([]*pbc.Element, len(allCoef))
	com := PolyCommit{
		c: C,
	}
	for i, coeff := range allCoef {
		com.c[i] = Curve.Pairing.NewG1()
		pow := conv.GmpInt2BigInt(coeff)
		Curve.Element.PowBig(com.c[i], pow)
	}
	return com
}
func (com PolyCommit) Verify(poly poly.Poly) bool {
	coeff := poly.GetAllCoeff()
	C := make([]*pbc.Element, len(coeff))
	commCheck := PolyCommit{
		c: C,
	}
	for i, tmp := range coeff {
		commCheck.c[i] = Curve.Pairing.NewG1()
		Curve.Element.PowBig(commCheck.c[i], conv.GmpInt2BigInt(tmp))
		if !commCheck.c[i].Equals(com.c[i]) {
			return false
		}
	}
	return true
}
func (com PolyCommit) VerifyEval(x *big.Int, y *big.Int) bool {
	gYRef := Curve.Pairing.NewG1()
	Curve.Element.PowBig(gYRef, y)

	xx := big.NewInt(1)

	gPx := Curve.Pairing.NewG1()
	gPx.Set1()

	tmp := Curve.Pairing.NewG1()
	for i := range com.c {
		// tmp = g^ai^{x^i}
		tmp.PowBig(com.c[i], xx)

		gPx.Mul(tmp, gPx)

		xx.Mul(xx, x)
		xx.Mod(xx, Curve.BigInt)
	}

	return gPx.Equals(gYRef)
}
func AdditiveHomomorphism(commQ, commR PolyCommit) PolyCommit {
	if len(commQ.c) != len(commR.c) {
		panic("mismatch degree")
	}

	comm := PolyCommit{
		c: make([]*pbc.Element, len(commQ.c)),
	}

	for i := range comm.c {
		comm.c[i] = Curve.Pairing.NewG1()
		comm.c[i].Mul(commQ.c[i], commR.c[i])
	}

	return comm
}
