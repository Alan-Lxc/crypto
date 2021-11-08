package poly

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPolynomial_Mul(t *testing.T) {
	op1 := FromVec(1, 1, 1, 1, 1, 1)
	result := NewEmpty()

	err := result.Multiply(op1, op1)
	assert.Nil(t, err, "Mul")
	//fmt.Println(result.GetAllCoeff())
	expected := FromVec(1, 2, 3, 4, 5, 6, 5, 4, 3, 2, 1)
	assert.True(t, expected.IsSame(result), "Mul")
}
