// Copyright 2013 com authors
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

import "regexp"

const (
	regex_email_pattern        = `(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`
	regex_strict_email_pattern = `(?i)[A-Z0-9!#$%&'*+/=?^_{|}~-]+` +
		`(?:\.[A-Z0-9!#$%&'*+/=?^_{|}~-]+)*` +
		`@(?:[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?\.)+` +
		`[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?`
	regex_url_pattern  = `(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?`
	regex_psw_number   = `^[0-9]*$`    // 验证密码是否为纯数字
	regex_psw_alphabet = `^[A-Za-z]+$` // 验证密码是否为纯英文字母
)

var (
	regex_email        *regexp.Regexp
	regex_strict_email *regexp.Regexp
	regex_url          *regexp.Regexp
	regex_psw_n        *regexp.Regexp
	regex_psw_s        *regexp.Regexp
)

func init() {
	regex_email = regexp.MustCompile(regex_email_pattern)
	regex_strict_email = regexp.MustCompile(regex_strict_email_pattern)
	regex_url = regexp.MustCompile(regex_url_pattern)
	regex_psw_n = regexp.MustCompile(regex_psw_number)
	regex_psw_s = regexp.MustCompile(regex_psw_alphabet)
}

// validate string is an email address, if not return false
// basically validation can match 99% cases
func IsEmail(email string) bool {
	return regex_email.MatchString(email)
}

// validate string is an email address, if not return false
// this validation omits RFC 2822
func IsEmailRFC(email string) bool {
	return regex_strict_email.MatchString(email)
}

// validate string is a url link, if not return false
// simple validation can match 99% cases
func IsUrl(url string) bool {
	return regex_url.MatchString(url)
}

// If password verification - number
func IsPswNum(str string) bool {
	return regex_psw_n.MatchString(str)
}

// If password verification - letters
func IsPswLetters(str string) bool {
	return regex_psw_s.MatchString(str)
}

// 组合验证密码是否为弱密码
func IsWeakPsw(str string) bool {
	// 判断密码是否为纯数字
	if regex_psw_n.MatchString(str) {
		return true
	}

	// 判断密码是否为纯字母
	if regex_psw_s.MatchString(str) {
		return true
	}
	return false
}
