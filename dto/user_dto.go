package dto

import "baseGo/model"

type UserDto struct {
	Name    string `json:"name"`
	Account string `json:"account"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name:    user.Name,
		Account: user.Account,
	}
}
