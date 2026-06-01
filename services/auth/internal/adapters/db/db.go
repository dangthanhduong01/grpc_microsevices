package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Adapter struct {
	db *gorm.DB
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"uniqueIndex;not null"`
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	FullName string `gorm:"default:""`
}

func NewAdapter(dataSourceURL string) (*Adapter, error) {
	db, err := gorm.Open(postgres.Open(dataSourceURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}

	return &Adapter{db: db}, nil
}

func (a *Adapter) CreateUser(email, username, password string) (string, error) {
	user := User{
		Email:    email,
		Username: username,
		Password: password, // In production, hash the password
	}

	result := a.db.Create(&user)
	if result.Error != nil {
		return "", result.Error
	}
	return fmt.Sprintf("%d", user.ID), nil
}

func (a *Adapter) GetUserByEmail(email string) (string, error) {
	var user User
	result := a.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return "", result.Error
	}
	return fmt.Sprintf("%d", user.ID), nil
}

func (a *Adapter) GetUserByUsername(username string) (string, error) {
	var user User
	result := a.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return "", result.Error
	}
	return fmt.Sprintf("%d", user.ID), nil
}

func (a *Adapter) GetUserByID(id string) (string, error) {
	var user User
	result := a.db.First(&user, id)
	if result.Error != nil {
		return "", result.Error
	}
	return fmt.Sprintf("%d", user.ID), nil
}

func (a *Adapter) UpdateUser(id, email, username, fullName string) error {
	var user User
	result := a.db.First(&user, id)
	if result.Error != nil {
		return result.Error
	}

	if email != "" {
		user.Email = email
	}
	if username != "" {
		user.Username = username
	}
	if fullName != "" {
		user.FullName = fullName
	}

	result = a.db.Save(&user)
	return result.Error
}