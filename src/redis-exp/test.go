package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func main() {
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

	bytes := [4]byte{1, 2, 3, 4}
	str := convert(bytes[:])

	fmt.Println(str, "str")

}

func convert(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, ",")
}
