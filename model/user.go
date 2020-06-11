package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"varchar(20);"`
	Account   string `gorm:"varchar(20);not null"`
	Telephone string `gorm:"varchar(11);"`
	Password  string `gorm:"varchar(255);not null"`
}
