package point

import (
	"github.com/Nik-U/pbc"
	"github.com/ncw/gmp"
)
type point struct {
	x	*gmp.Int
	y	*gmp.Int
	PolyWit	*pbc.Element
}
func newPoint(x *gmp.Int,y *gmp.Int,w *pbc.Element)*point{
return	&point{
	x: x,
	y: y,
	PolyWit: w,
}
}
