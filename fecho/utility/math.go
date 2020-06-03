// Copyright 2014 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package utility

import (
	"fmt"
	"math"
	"strconv"
)

const (
	ACCURACY = 0.00001 //精度
)

// PowInt is int type of math.Pow function.
func PowInt(x int, y int) int {
	num := 1
	for i := 0; i < y; i++ {
		num *= x
	}
	return num
}

// ToFixed10K 5位小数精度
func ToFixed10K(l float64) float64 {
	if math.IsNaN(l) {
		return 0
	}

	//return float64(int(l*100000)) / 100000
	s := fmt.Sprintf("%.5f", l)
	inst, _ := strconv.ParseFloat(s, 64)
	return inst
}

// ToFixed 自定义小数位精度.
func ToFixed(l float64, n float64) float64 {
	if math.IsNaN(l) {
		return 0
	}
	if l < 0 {
		return float64(int((l-0.5/n)*n)) / n
	}

	return float64(int((l+0.5/n)*n)) / n
}

// 保留两位小数(不做四舍五入)
func RoundTwo(l, n float64) float64 {
	if math.IsNaN(l) {
		return 0
	}

	return float64(int(l*n)) / n
}

// If 三元表达式
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// BeBig 比大
func BeBig(source, compare float64) bool {
	if source > compare || compare-source < ACCURACY {
		return true
	}
	return false
}

// BeBigOrEqual 大于等于
func BeBigOrEqual(source, compare float64) bool {
	if source > compare || math.Abs(compare-source) <= ACCURACY {
		return true
	}
	return false
}

// BeSmall 比小
func BeSmall(source, compare float64) bool {
	if source < compare || source-compare < ACCURACY {
		return true
	}
	return false
}

// BeSmallOrEqual 小于等于
func BeSmallOrEqual(source, compare float64) bool {
	if source < compare || math.Abs(source-compare) < ACCURACY {
		return true
	}
	return false
}

// 相等
func BeEqual(source, compare float64) bool {
	if math.Abs(source-compare) < ACCURACY {
		return true
	}
	return false
}

// Round 浮点型精准问题矫正
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

type Accuracy func() float64

func (this Accuracy) Equal(a, b float64) bool {
	return math.Abs(a-b) < this()
}

func (this Accuracy) Greater(a, b float64) bool {
	return math.Max(a, b) == a && math.Abs(a-b) > this()
}

func (this Accuracy) GreaterOrEqual(a, b float64) bool {
	return math.Max(a, b) == a || math.Abs(a-b) < this()
}

func (this Accuracy) Smaller(a, b float64) bool {
	return math.Max(a, b) == b && math.Abs(a-b) > this()
}

func (this Accuracy) SmallerOrEqual(a, b float64) bool {
	return math.Max(a, b) == b || math.Abs(a-b) < this()
}

var AccuracyComp Accuracy = func() float64 { return ACCURACY }

func ToFixedBy2DecimalPlaces(l float64) float64 {
	s := fmt.Sprintf("%.2f", l)
	inst, _ := strconv.ParseFloat(s, 64)
	return inst
}

func ToFixedByNDecimalPlaces(l float64, precision int) float64 {
	if precision < 0 {
		precision = 0
	}
	fmtByPrecision := fmt.Sprintf(fmt.Sprintf("%%.%df", precision), l)
	inst, _ := strconv.ParseFloat(fmtByPrecision, 64)
	return inst
}
