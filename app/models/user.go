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

// Callback before update user
func (m *User) BeforeUpdate() (err error) {
    m.UpdatedAt = time.Now()
    return
}

// Callback before create user
func (m *User) BeforeCreate() (err error) {
    m.CreatedAt = time.Now()
    return
}

// Create
func Create(m *User) (*User, error) {
	var err error
	err = orm.Create(&m)
	return m, err
}

// Update
func (m *User) Update() error {
	var err error
	err = orm.Save(&m)
	return err
}

// Delete
func (m *User) Delete() error {
	var err error
	err = orm.Delete(&m)
	return err
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
	err = orm.FindAll(users)
	return users, err
}
