package common

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"log"
	"math/rand"
	"model"
	"model/code"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type H map[string]interface{}

//判断是否为数字
func IsNumeric(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
	case float32, float64, complex64, complex128:
		return true
	case string:
		str := val.(string)
		if str == "" {
			return false
		}
		// Trim any whitespace
		str = strings.Trim(str, " \\t\\n\\r\\v\\f")
		if str[0] == '-' || str[0] == '+' {
			if len(str) == 1 {
				return false
			}
			str = str[1:]
		}
		// hex
		if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
			for _, h := range str[2:] {
				if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
					return false
				}
			}
			return true
		}
		// 0-9,Point,Scientific
		p, s, l := 0, 0, len(str)
		for i, v := range str {
			if v == '.' { // Point
				if p > 0 || s > 0 || i+1 == l {
					return false
				}
				p = i
			} else if v == 'e' || v == 'E' { // Scientific
				if i == 0 || s > 0 || i+1 == l {
					return false
				}
				s = i
			} else if v < '0' || v > '9' {
				return false
			}
		}
		return true
	}

	return false
}

// 判断obj是否在target中，target支持的类型arrary,slice,map
func Contain(target interface{}, obj interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		by, _ := json.Marshal(obj)
		h := md5.New()
		h.Write(by)
		objmd5 := hex.EncodeToString(h.Sum(nil))
		for i := 0; i < targetValue.Len(); i++ {
			byt, _ := json.Marshal(targetValue.Index(i).Interface())
			ht := md5.New()
			ht.Write(byt)
			tgtmd5 := hex.EncodeToString(ht.Sum(nil))
			if objmd5 == tgtmd5 {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func IsNotBlank(str string) bool {

	if len(str) == 0 {
		return false
	}

	if strings.TrimSpace(str) == "" {
		return false
	}

	return true
}

func IsBlank(str string) bool {

	if len(str) == 0 {
		return true
	}

	if len(strings.TrimSpace(str)) == 0 {
		return true
	}

	return false
}

func CountString(msg string) int64 {

	var total int64

	for range msg {
		total++
	}

	return total
}

//过滤空白字符
func ReplaceBlank(str string) string {

	if len(str) == 0 {
		return str
	}

	charactors := make([]rune, 0)

	for _, v := range str {
		if unicode.IsSpace(v) {
			charactors = append(charactors, ' ')
		} else {
			charactors = append(charactors, v)
		}
	}

	return string(charactors)
}

func CompileTemplate(template, variable, sign string) string {

	if IsBlank(template) {
		return ""
	}

	template = ReplaceBlank(template) //去掉\t\n等特殊字符

	if IsBlank(variable) {
		return template + "【" + sign + "】"
	}

	variable = strings.TrimSpace(variable)

	variable = ReplaceBlank(variable) //去掉\t\n等特殊字符

	jo := make(map[string]interface{})

	err := json.Unmarshal([]byte(variable), &jo)
	if err != nil {
		log.Println("json.Unmarshal error: variable ==>", variable)
		fmt.Println("json.Unmarshal error: variable ==>", variable)
		return template + "【" + sign + "】"
	}

	for k, v := range jo {

		var str string

		switch v.(type) {
		case string:
			str = fmt.Sprintf("%s", v)
		case int:
			str = fmt.Sprintf("%d", v)
		case int32:
			str = fmt.Sprintf("%d", v)
		case int64:
			str = fmt.Sprintf("%ld", v)
		case float64:
			str = fmt.Sprintf("%.2f", v)
		case float32:
			str = fmt.Sprintf("%.2f", v)
		}

		template = strings.Replace(template, "{"+k+"}", str, -1)
	}

	return template + "【" + sign + "】"
}

// Return part of a string
//
// Examples:
// string   start length return
// "abcdef" 0     1      a
// "abcdef" 1     2      bc
// "abcdef" -2    1      e
// "abcdef" -2    0      ef
// "abcdef" 1     -2     bcd
// "abcdef" -3    -2     d
// "abcdef" -2    -4     ""
// "abcdef" -20   -4     ab
// "abcdef" -20   -10    ""
func SubStr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl + start
		if start < 0 {
			start = 0
		}
	}

	if 0 == length {
		end = rl
	} else if 0 > length {
		end = rl + length
		if end < 0 {
			end = 0
		}
	} else {
		end = start + length
		if end > rl {
			end = rl
		}
	}

	if start > end {
		return ""
	}

	return string(rs[start:end])
}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {

	str := time.Time(t).Format("2006-01-02 15:04:05")

	return []byte(str), nil
}
func Truncation(t float64) float64 {
	t, _ = strconv.ParseFloat(strconv.FormatFloat(t, 'f', 2, 64), 64)
	return t
}
func RandomInt() int {
	return rand.Intn(100)
}
func RandomFloat() float64 {
	return rand.Float64()
}
func PasswordCheck(pwd string) (bool, int, error) {
	//判断是否在A-Z
	bl, err := regexp.MatchString("[A-Z]+", pwd)
	if err != nil {
		return false, 0, err

	}
	if !bl {
		return false, code.NO_UPPERCASE_LETTER, err
	}

	//判断是否在a-z
	bl, err = regexp.MatchString("[a-z]+", pwd)
	if err != nil {
		return false, 0, err

	}
	if !bl {
		return false, code.NO_LOWERCASE_LETTER, err
	}

	//判断是否在1-9
	bl, err = regexp.MatchString("[0-9]+", pwd)
	if err != nil {
		return false, 0, err

	}
	if !bl {
		return false, code.NO_FIGURES, err
	}

	//判断是否有非法字符
	bl, err = regexp.MatchString("([^a-z0-9A-Z])+", pwd)
	if err != nil {
		return false, 0, err

	}
	if bl {
		return false, code.ILLEGAL_CHARACTERS, err
	}
	return true, 0, nil
}

//decimal
func DecimalSum(a float64, b float64) float64 {
	af := decimal.NewFromFloat(a)
	bf := decimal.NewFromFloat(b)
	sf, _ := decimal.Sum(af, bf).Round(model.ROUND_TWO).Float64()
	return sf
}

//
func DecimalMul(a float64, b float64) float64 {
	af := decimal.NewFromFloat(a)
	bf := decimal.NewFromFloat(b)
	sf, _ := af.Mul(bf).Round(model.ROUND_TWO).Float64()
	return sf
}

//
func DecimalSub(a float64, b float64) float64 {
	af := decimal.NewFromFloat(a)
	bf := decimal.NewFromFloat(b)
	sf, _ := af.Sub(bf).Round(model.ROUND_TWO).Float64()
	return sf
}

//
func DecimalDiv(a float64, b float64) float64 {
	if b == 0.0 {
		return 0.0
	}
	af := decimal.NewFromFloat(a)
	bf := decimal.NewFromFloat(b)
	sf, _ := af.Div(bf).Round(model.ROUND_TWO).Float64()
	return sf
}

// 保留五位小数
func DecimalAddByRoundFive(a float64, b float64) float64 {
	af := decimal.NewFromFloat(a)
	bf := decimal.NewFromFloat(b)
	sf, _ := af.Add(bf).Round(model.ROUND_FIVE).Float64()
	return sf
}

// 保留五位小数
func DecimalMulByRoundFive(a float64, b float64) float64 {
	af := decimal.NewFromFloat(a)
	bf := decimal.NewFromFloat(b)
	sf, _ := af.Mul(bf).Round(model.ROUND_FIVE).Float64()
	return sf
}

// 保留五位小数
func DecimalSubByRoundFive(a float64, b float64) float64 {
	af := decimal.NewFromFloat(a)
	bf := decimal.NewFromFloat(b)
	sf, _ := af.Sub(bf).Round(model.ROUND_FIVE).Float64()
	return sf
}

// 保留五位小数
func DecimalDivByRoundFive(a float64, b float64) float64 {
	af := decimal.NewFromFloat(a)
	bf := decimal.NewFromFloat(b)
	sf, _ := af.Div(bf).Round(model.ROUND_FIVE).Float64()
	return sf
}

// 不处理小数
func DecimalAdds(a float64, b float64) float64 {
	af := decimal.NewFromFloat(a)
	bf := decimal.NewFromFloat(b)
	sf, _ := af.Add(bf).Float64()
	return sf
}

// 不处理小数
func DecimalMuls(a float64, b float64) float64 {
	af := decimal.NewFromFloat(a)
	bf := decimal.NewFromFloat(b)
	sf, _ := af.Mul(bf).Float64()
	return sf
}

// 不处理小数
func DecimalSubs(a float64, b float64) float64 {
	af := decimal.NewFromFloat(a)
	bf := decimal.NewFromFloat(b)
	sf, _ := af.Sub(bf).Float64()
	return sf
}

// 不处理小数
func DecimalDivs(a float64, b float64) float64 {
	af := decimal.NewFromFloat(a)
	bf := decimal.NewFromFloat(b)
	sf, _ := af.Div(bf).Float64()
	return sf
}
