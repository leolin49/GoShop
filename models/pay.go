package models

import (
	"time"

	"gorm.io/gorm"
)

type PaymentLog struct {
	gorm.Model
	UserId	uint32
	OrderId string
	TransactionId string
	Amount	float32
	PayAt	time.Time
}

func (p PaymentLog) TableName() string {
	return "payment_log"
}

type PaymentLogQuery struct {
	db *gorm.DB
}

func NewPaymentLogQuery(db *gorm.DB) *PaymentLogQuery {
	return &PaymentLogQuery{
		db: db,
	}
}

func (q *PaymentLogQuery) CreatePaymentLog(payment *PaymentLog) error {
	return q.db.Model(&PaymentLog{}).Create(payment).Error
}
