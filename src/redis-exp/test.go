package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type xv struct{}

func main() {
	demos()

	bytes := [4]byte{1, 2, 3, 4}
	str := convert(bytes[:])

	fmt.Println(str, "str")

	var x interface{} = func(i int) float64 {
		i = 100
		return float64(i)
	}
	tryTypes(x)

}

func (*xv) tryTypes(x) {

	switch i := x.(type) {
	case nil:
		fmt.Println("x is nil")
	case int:
		fmt.Println(i)
	case func(int) float64:
		fmt.Println("type is func")
	default:
		fmt.Println("don`t know the type")
	}

}

func convert(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, ",")
}

func demos() {
	b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)

	var f interface{}
	json.Unmarshal(b, &f)

	m := f.(map[string]interface{})
	fmt.Println(m["Parents"])  // 读取 json 内容
	fmt.Println(m["a"] == nil) // 判断键是否存在

	var src = []int{1, 2, 3, 4}
	var temp = make([]string, len(src))
	for k, v := range src {
		temp[k] = fmt.Sprintf("%d", v)
	}
	var result = "[" + strings.Join(temp, ",") + "]"
	fmt.Println(result, src)
}

// func (*t) functionOfSomeType() error {
// 	// switch t := t.(type) {
// 	// default:
// 	// 	fmt.Printf("unexpected type %T", t) // %T prints whatever type t has
// 	// case bool:
// 	// 	fmt.Printf("boolean %t\n", t) // t has type bool
// 	// case int:
// 	// 	fmt.Printf("integer %d\n", t) // t has type int
// 	// case *bool:
// 	// 	fmt.Printf("pointer to boolean %t\n", t) // t has type bool
// 	// case *int:
// 	// 	fmt.Printf("pointer to integer %d\n", t) // t has type int
// 	// }
// 	// return nil

//  var t interface{}
// 	switch i := t.(type) {
// 	case nil:
// 		fmt.Println("x is nil") // i的类型是 x的类型 (interface{})
// 	case int:
// 		fmt.Println(i) // i的类型 int
// 	case float64:
// 		fmt.Println(i) // i的类型是 float64
// 	case func(int) float64:
// 		fmt.Println(i) // i的类型是 func(int) float64
// 	case bool, string:
// 		fmt.Println("type is bool or string") // i的类型是 x (interface{})
// 	default:
// 		fmt.Println("don't know the type") // i的类型是 x的类型 (interface{})
// 	}
// 	return nil
// }
