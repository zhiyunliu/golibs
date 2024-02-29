package xtypes

import (
	"github.com/shopspring/decimal"
)

type Decimal = decimal.Decimal
type NullDecimal = decimal.NullDecimal

var DecimalZero = decimal.Zero

var NewDecimal = decimal.New
var MaxDecimal = decimal.Max
var MinDecimal = decimal.Min
var AvgDecimal = decimal.Avg
var SumDecimal = decimal.Sum

var NewDecimalFromInt = decimal.NewFromInt
var NewDecimalFromInt32 = decimal.NewFromInt32
var NewDecimalFromBigInt = decimal.NewFromBigInt
var NewDecimalFromString = decimal.NewFromString
var NewDecimalFromFloat = decimal.NewFromFloat
var NewDecimalFromFloat32 = decimal.NewFromFloat32
var NewDecimalFromFloatWithExponent = decimal.NewFromFloatWithExponent
