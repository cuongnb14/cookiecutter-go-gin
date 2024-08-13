package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Firstname string `gorm:"size:255;not null" json:"firstname"`
	Lastname  string `gorm:"size:255;not null" json:"lastname"`
	Email     string `gorm:"size:100;not null;unique" json:"email"`
	Password  string `gorm:"size:100;not null;" json:"password"`
}
