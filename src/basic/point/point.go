package point

import (
	"github.com/Nik-U/pbc"
	"github.com/ncw/gmp"
)

type Point struct {
	X       *gmp.Int
	Y       *gmp.Int
	PolyWit *pbc.Element
}

func NewPoint(x *gmp.Int, y *gmp.Int, w *pbc.Element) *Point {
	return &Point{
		X:       x,
		Y:       y,
		PolyWit: w,
	}
}

type Pointmsg struct {
	index int
	point *Point
}

func (pointmsg *Pointmsg) GetIndex() int {
	if pointmsg != nil {
		return pointmsg.index
	} else {
		return 0
	}
}
func (pointmsg *Pointmsg) GetPoint() *Point {
	if pointmsg != nil {
		return pointmsg.point
	} else {
		return nil
	}
}
