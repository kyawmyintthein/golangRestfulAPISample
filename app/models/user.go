package models

import (
	"golangRestfulAPISample/app/models/orm"
	"time"
)

type (
	User struct {
		BaseModel
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

// Create
func Create(m *User) (User, error) {
	var err error
	err = orm.Create(&m)
	return m, err
}

// Update
func (m *User) Update() error {
	var err error
	err = orm.Save(&m)
	return m, err
}

// Delete
func (m *User) Delete() error {
	var err error
	err = orm.Delete(&m)
	return m, err
}

// FindUserByID
func FindUserByID(id uint64) (User, error) {
	var (
		user User
		err  error
	)
	err = orm.FindOneByID(&user, id)
	return user, err
}

// FindAllUsers
func FindAllUsers() ([]User, error) {
	var (
		users []User
		err   error
	)
	err = orm.FindOneByID(&users)
	return users, err
}
