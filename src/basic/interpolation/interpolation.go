package interpolation

import (
	"errors"
	. "github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	"github.com/ncw/gmp"
)

//Get a polynomial that satisfy all x and y
//degree
//func LagrangeInterpolate(degree int, x []*gmp.Int, y []*gmp.Int, mod *gmp.Int) (Poly, error) {
//
//	//使用多项式l_i来存储中间多项式l_i(x),l_i[i]=l_i(i)
//	//l(x) = (x-1)(x-2)...(x-n)
//	//l_i(x) = l(x)/(x-i)
//	//lambda_i(x) = l_i(x) * y[i] / l_i(i)
//	//这里令l(x)为product，即全排列；令l_i(x)为numerator，即分子，令l_i(i)为denominator, 即分母
//	//product初始化为只有常数1的多项式
//
//	//初始化变量
//
//	//tmp为临时多项式，初始化为默认一次多项式
//	tmpPoly, err := NewPoly(1)
//	numerator, err := NewPoly(degree)
//	permutation := NewConstant(1)
//	//resultPoly为结果多项式，即拉格朗日插值多项式
//	resultPoly, err := NewPoly(degree)
//
//	if err != nil {
//		return Poly{}, err
//	}
//	//denominator初始化为Int0
//	denominator := gmp.NewInt(0)
//
//	//求得product
//	//首先置tmp一次项为1
//	tmpPoly.SetCoeffWithInt(1, 1)
//	tmpInt, err := tmpPoly.GetCoeff(0)
//	for i := 0; i <= degree; i++ {
//		tmpInt, err = tmpPoly.GetCoeff(0)
//		tmpPoly.SetCoeffWithGmp(0, tmpInt.Neg(x[i]))
//		permutation.Multiply(permutation, tmpPoly)
//	}
//
//	//依此求得拉格朗日分式，并相加，注意要模mod
//	for i := 0; i <= degree; i++ {
//		//每次循环开始时把分母置1
//		denominator.Set(gmp.NewInt(1))
//
//		//计算分母多项式
//		tmpInt, err = tmpPoly.GetCoeff(0)
//		tmpPoly.SetCoeffWithGmp(0, tmpInt.Neg(x[i]))
//
//		err = numerator.Divide(permutation, tmpPoly)
//		if err != nil {
//			return Poly{}, err
//		}
//		numerator.Mod(mod)
//		//使用分母多项式带入计算出分母真实值
//
//		//这样计算出来得到分母的真实值
//		numerator.EvalMod(x[i], mod, denominator)
//		//检测分母真实值是否等于0，一般不会等于0，如果等于0可能是有内鬼
//		if 0 == denominator.CmpInt32(0) {
//			return Poly{}, errors.New("internal error: check dupliction in x[]")
//		}
//
//		//通过分母真实值求分子真实值，先求分母真实值的模逆
//		denominator.ModInverse(denominator, mod)
//		denominator.Mul(denominator, y[i])
//		resultPoly.AddMul(numerator, denominator)
//
//	}
//	//最后再将结果取模
//	resultPoly.Mod(mod)
//	return resultPoly, nil
//
//}
// LagrangeInterpolate returns a polynomial of specified degree that pass through all points in x and y
func LagrangeInterpolate(degree int, x []*gmp.Int, y []*gmp.Int, mod *gmp.Int) (Poly, error) {
	// initialize variables
	tmp, err := NewPoly(1)
	if err != nil {
		return Poly{}, err
	}

	inter, err := NewPoly(degree)
	if err != nil {
		return Poly{}, err
	}

	product := NewConstant(1)

	resultPoly, err := NewPoly(degree)
	if err != nil {
		return Poly{}, err
	}
	//分母=0，新对象
	denominator := gmp.NewInt(0)
	//
	// tmp(x) = x - x[i]
	//置一次项系数 为1
	tmp.SetCoeffWithInt(1, 1)
	// note only the first degree points are used
	//得到全排列，product就是l(x),每一位代表一个多项式系数
	for i := 0; i <= degree; i++ {
		tmp.GetPtrtoConstant().Neg(x[i])
		//product反复乘以tmp
		product.MulSelf(tmp)
	}

	for i := 0; i <= degree; i++ {
		denominator.Set(gmp.NewInt(1))
		// compute denominator and numerator

		// tmp = x - x[i]
		tmp.SetCoeffWithInt(1, 1) // i don't think this needed...
		tmp.GetPtrtoConstant().Neg(x[i])

		// inner(x) = (x-1)(x_2)...(x-n) except for (x-i)
		//product 是全排列，inter是分子
		err = inter.Divide(product, tmp)
		if err != nil {
			return Poly{}, err
		}

		// lambda_i(x) = inner(x) * y[i] / inner(x[i])

		//分子多项式取模，把分子带入求解得到分母。利用分子得到分母
		inter.Mod(mod)
		inter.EvalMod(x[i], mod, denominator)

		//如果分母为0报错，一般不会为0
		// panic if denominator == 0
		if 0 == denominator.CmpInt32(0) {
			return Poly{}, errors.New("internal error: check duplication in x[]")
		}

		//inter 是分子，
		//denominator 是分母，
		//y[i]是函数值
		//求分母模逆
		denominator.ModInverse(denominator, mod)
		//分母真实值乘以y[i]就得到分母的真实值
		denominator.Mul(denominator, y[i])
		//最后结果加上inter和分母模逆的乘积
		resultPoly.AddMul(inter, denominator)

	}

	resultPoly.Mod(mod)

	return resultPoly, nil
}
