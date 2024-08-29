package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex"`
	Name     string
	GoogleID string `gorm:"uniqueIndex"`
}

func UpsertUser(db *gorm.DB, email, name, googleID string) (*User, error) {
	var user User
	result := db.Where(User{GoogleID: googleID}).Assign(User{Email: email, Name: name}).FirstOrCreate(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
