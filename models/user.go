package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type User struct {
	gorm.Model
	Email       string
	Name        string
	Password    string
	Age         uint8
	Birthday    *time.Time
	PhoneNumber *string // A pointer to a string, allowing for null values.
	Address     *string
}

func (u User) TableName() string {
	return "user"
}

type UserQuery struct {
	db *gorm.DB
}

func NewUserQuery(db *gorm.DB) *UserQuery {
	return &UserQuery{
		db: db,
	}
}

func NewUserQueryWrite(db *gorm.DB) *UserQuery {
	return &UserQuery{
		db: db.Clauses(dbresolver.Write),
	}
}

func NewUserQueryRead(db *gorm.DB) *UserQuery {
	return &UserQuery{
		db: db.Clauses(dbresolver.Write),
	}
}

func NewUserQueryWithDBName(db *gorm.DB, dbName string) *UserQuery {
	return &UserQuery{
		db: db.Clauses(dbresolver.Use(dbName)),
	}
}

func (q *UserQuery) GetUserIdByEmail(email string) (uint32, error) {
	var user User
	err := q.db.Model(User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return uint32(user.ID), nil
}

func (q *UserQuery) GetIdAndPwdByEmail(email string) (uint32, string, error) {
	var user User
	err := q.db.Select("id, password").Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, "", errors.New("user not exist")
		} else {
			return 0, "", err
		}
	}
	return uint32(user.ID), user.Password, nil
}

func (q *UserQuery) CreateUser(user *User) (uint32, error) {
	// check again
	if user_id, err := q.GetUserIdByEmail(user.Email); err != nil {
		return 0, err
	} else if user_id == 0 {
		return 0, errors.New(
			fmt.Sprintf("create user failed cause user [%s] already exist", user.Email),
		)
	}
	// create
	err := q.db.Create(user).Error
	if err != nil {
		return 0, err
	}
	return q.GetUserIdByEmail(user.Email)
}
