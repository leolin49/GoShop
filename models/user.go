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

func (u *User) TableName() string {
	return "user"
}

// NOTE: Hook function
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return nil
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

func (q *UserQuery) GetUserIdByEmail(email string) (userId uint32, err error) {
	err = q.db.Model(&User{}).Select("id").Where("email = ?", email).First(&userId).Error
	return
}

func (q *UserQuery) CreateUser(user *User) error {
	return q.db.Transaction(func(tx *gorm.DB) error {
		// check
		var checkUserId uint32
		if err := tx.Model(&User{}).Select("id").Where("email = ?", user.Email).First(&checkUserId).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		} else if checkUserId > 0 {
			return errors.New(
				fmt.Sprintf("create user [%d] failed cause user [%v] already exist", checkUserId, user),
			)
		}
		// create
		if err := tx.Create(user).Error; err != nil {
			tx.Rollback()
			return tx.Error
		}

		// commit transaction
		return nil
	})
}
