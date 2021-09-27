package nodes

import (
	"github.com/ncw/gmp"
	"math/rand"
)

type Node struct {
	//label of Node
	label int
	//Total number of Nodes
	total int
	//degree of polynomial
	degree int
	//the polynomial was set on Z_p
	p *gmp.Int

	// Rand source
	randstate *rand.Rand
}
