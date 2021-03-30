package common

import "github.com/awatercolorpen/godruid"

func Dimension(dimension string) interface{} {
	return godruid.DimDefault(dimension, dimension)
}
