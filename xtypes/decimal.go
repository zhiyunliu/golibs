package xtypes

import (
	"github.com/shopspring/decimal"
)

type Decimal = decimal.Decimal

var NewDecimal = decimal.New

var NewDecimalFromInt = decimal.NewFromInt
var NewDecimalFromInt32 = decimal.NewFromInt32
var NewDecimalFromBigInt = decimal.NewFromBigInt
var NewDecimalFromString = decimal.NewFromString
var NewDecimalFromFloat = decimal.NewFromFloat
var NewDecimalFromFloat32 = decimal.NewFromFloat32
var NewDecimalFromFloatWithExponent = decimal.NewFromFloatWithExponent
