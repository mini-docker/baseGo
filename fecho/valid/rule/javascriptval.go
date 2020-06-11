package rule

import (
	"strings"
)

// 是否含有javascript
func JavascriptFilter(src string) bool {

	if src == "" {
		return false
	}

	// 处理SCRIPT
	hasScriptTag := false
	hasScriptTag = strings.Contains(src, "script")
	//处理%
	hasPaTag := false
	hasPaTag = strings.Contains(src, "%")
	//处理_
	hasUnag := false
	hasUnag = strings.Contains(src, "%")

	return hasScriptTag || hasPaTag || hasUnag
}
