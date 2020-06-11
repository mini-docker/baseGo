/*
 * Copyright 2018 The golyu Authors. All rights reserved.
 * Use of this source code is governed by a Apache License, Version 2.0 (the "License");
 * license that can be found in the LICENSE file.
 */

package rule

import (
	"errors"
	//"im_chat/js108/fecho/echo"
	"regexp"
)

const (
	NO_UPPERCASE_LETTER      = 8606
	NO_LOWERCASE_LETTER      = 8607
	NO_FIGURES               = 8608
	ILLEGAL_CHARACTERS       = 8609
	ILLEGAL_CHARACTERS_ERROR = 11120
)

// 校验password
type PasswordRule struct {
	value interface{}
	FullTag
}

func (r *PasswordRule) Clone() Rule {
	clone := *r
	return &clone
}

func (r *PasswordRule) Tag() string {
	return "Password"
}

func (r *PasswordRule) Generate(value interface{}, tagValue string) error {
	if value == nil {
		return errors.New("Generate Password:value is nil")
	}
	r.value = value
	if tagValue != r.Tag() {
		return errors.New("Generate Password:the tag are out of specification")
	}

	return nil
}

func (r *PasswordRule) Valid() error {
	if r.value == nil {
		return errors.New("Validation Password:value is nil")
	}
	switch r.value.(type) {

	case string:
		value := r.value.(string)
		bl, _, err := PasswordCheck(value)
		if err != nil {
			return err
		}
		if !bl {
			return errors.New("password illegal")
		}
	}
	return nil
}

func PasswordCheck(pwd string) (bool, int, error) {
	//判断是否在A-Z
	bl, err := regexp.MatchString("[A-Z]+", pwd)
	if err != nil {
		return false, 0, err

	}
	if !bl {
		return false, NO_UPPERCASE_LETTER, err
	}

	//判断是否在a-z
	bl, err = regexp.MatchString("[a-z]+", pwd)
	if err != nil {
		return false, 0, err

	}
	if !bl {
		return false, NO_LOWERCASE_LETTER, err
	}

	//判断是否在1-9
	bl, err = regexp.MatchString("[0-9]+", pwd)
	if err != nil {
		return false, 0, err

	}
	if !bl {
		return false, NO_FIGURES, err
	}

	//判断是否有非法字符
	bl1, err := regexp.MatchString(`[\\]`, pwd)
	if err != nil {
		return false, 0, err

	}

	bl2, err := regexp.MatchString("[-~!@#$^&*()_+=|{}'\":,.\\[\\]<>/?！%￥…（）【】‘：”“。，、？·「」]", pwd)
	if err != nil {
		return false, 0, err

	}
	if !bl1 && !bl2 {
		return false, ILLEGAL_CHARACTERS_ERROR, err
	}
	return true, 0, nil
}

// 取款密码校验
func DrawPasswordCheck(pwd string) bool {
	//判断是否在A-Z
	flag, err := regexp.MatchString("([0-9]){6}", pwd)
	if err != nil {
		return false

	}
	if !flag {
		return false
	}
	return true
}
