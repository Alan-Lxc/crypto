package point

import (
	"github.com/Nik-U/pbc"
	"github.com/ncw/gmp"
)

type Point struct {
	X       int32
	Y       *gmp.Int
	PolyWit *pbc.Element
}

func NewPoint(x int32, y *gmp.Int) *Point {
	return &Point{
		X: x,
		Y: y,
		//PolyWit: w,
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
func (pointmsg *Pointmsg) SetPoint(point *Point) {
	if point != nil {
		pointmsg.point = point
	}
}

//transport a object value to a method must use pointer
func (pointmsg *Pointmsg) SetIndex(index int) {
	if index > 0 {
		pointmsg.index = index
	}
}
