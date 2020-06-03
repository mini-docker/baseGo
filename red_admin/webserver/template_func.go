package webserver

import (
	"html/template"
	"strings"
)

func unescape(x string) interface{} {
	return template.HTML(x)
}

func multi(x int, y float32) float32 {
	return float32(x) * y
}

func upper(str string) string {
	return strings.ToUpper(str)
}

func mod(x, y, z int) bool {
	return x%y == z
}

// 通过id 获取名称
func category(id int) string {
	switch id {
	case 1:
		return "lottery"
	case 2:
		return "electronics"
	case 3:
		return "video"
	case 4:
		return "chess"
	case 5:
		return "sports"
	case 6:
		return "fish"
	default:
	}

	return ""
}
