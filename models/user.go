package models

import(
    "echo_rest_api/db/gorm"
    "time"
)

type (
	User struct {
        BaseModel
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

func CreateUser(m *User) error{
	var (
		err  error
	)
    m.CreatedAt = time.Now()
    tx := gorm.MysqlConn().Begin()
    if err = tx.Create(&m).Error; err != nil {
        tx.Rollback()
        return err
    }
    tx.Commit()
	return err
}

func (m *User) UpdateUser(data *User) error{
    var (
        err  error
    )

    m.UpdatedAt = time.Now()
    m.Name = data.Name
    m.Email = data.Email

    tx := gorm.MysqlConn().Begin()
    if err = tx.Save(&m).Error; err != nil {
        tx.Rollback()
        return err
    }
    tx.Commit()
    return err
}

func (m User) DeleteUser() error{
    var err error
    tx := gorm.MysqlConn().Begin()
    if err = tx.Delete(&m).Error; err != nil {
        tx.Rollback()
        return err
    }
    tx.Commit()
    return err
}

func GetUserById(id uint64) (User, error) {
    var (
        user User
        err  error
    )

    tx := gorm.MysqlConn().Begin()
    if err = tx.Last(&user, id).Error; err != nil {
        tx.Rollback()
        return user, err
    }
    tx.Commit()
    return user, err
}


func GetUsers() ([]User, error) {
    var (
        users []User
        err  error
    )

    tx := gorm.MysqlConn().Begin()
    if err = tx.Find(&users).Error; err != nil {
        tx.Rollback()
        return users, err
    }
    tx.Commit()
    return users, err
}


