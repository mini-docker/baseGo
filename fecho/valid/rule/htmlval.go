/*
 * Copyright 2018 The golyu Authors. All rights reserved.
 * Use of this source code is governed by a Apache License, Version 2.0 (the "License");
 * license that can be found in the LICENSE file.
 */

package rule

import (
	"errors"
	"html"
)

// 校验html
type HtmlRule struct {
	value interface{}
	FullTag
}

func (r *HtmlRule) Clone() Rule {
	clone := *r
	return &clone
}

func (r *HtmlRule) Tag() string {
	return "Html"
}

func (r *HtmlRule) Generate(value interface{}, tagValue string) error {
	if value == nil {
		return errors.New("Generate Html:value is nil")
	}
	r.value = value
	if tagValue != r.Tag() {
		return errors.New("Generate Html:the tag are out of specification")
	}

	return nil
}

func (r *HtmlRule) Valid() error {
	if r.value == nil {
		return errors.New("Validation Html:value is nil")
	}
	switch r.value.(type) {

	case string:
		value := r.value.(string)
		r.value = html.EscapeString(value)
	}
	return nil
}

// 是否含有HTML标签
func HtmlFilter(src string, clean bool, replace string) (string, bool, error) {

	if src == "" {
		return src, false, nil
	}

	// 转义标签
	return html.EscapeString(src), false, nil

	//
	//// 去除连续的换行符
	//re, err := regexp.Compile("\\s{2,}")
	//if err != nil {
	//	return "", false, err
	//} else {
	//	src = re.ReplaceAllString(src, "")
	//}
	//
	//// 将HTML标签全转换成小写
	//re, toLowerErr := regexp.Compile("\\<[\\S\\s]+?\\>")
	//if toLowerErr != nil {
	//	return "", false, toLowerErr
	//} else {
	//	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//}
	//
	//// 处理STYLE标签
	//b := []byte(src)
	//hasStyleTag := false
	//re, styleErr := regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	//if styleErr != nil {
	//	return "", false, styleErr
	//} else {
	//	if re.Match(b) {
	//		hasStyleTag = true
	//		if clean {
	//			src = re.ReplaceAllString(src, replace)
	//		}
	//	}
	//}
	//
	//// 处理SCRIPT标签
	//hasScriptTag := false
	//re, scriptErr := regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	//if scriptErr != nil {
	//	return "", false, scriptErr
	//} else {
	//	if re.Match(b) {
	//		hasScriptTag = true
	//		if clean {
	//			src = re.ReplaceAllString(src, replace)
	//		}
	//	}
	//}
	//
	////去除所有尖括号内的HTML代码，并换成换行符
	//hasHtmlTag := false
	//re, htmlErr := regexp.Compile("\\<[\\S\\s]+?\\>")
	//if htmlErr != nil {
	//	return "", false, htmlErr
	//} else {
	//	if re.Match(b) {
	//		hasHtmlTag = true
	//		if clean {
	//			src = re.ReplaceAllString(src, replace)
	//		}
	//	}
	//}

	return src, false, nil
}
