package models

import "gorm.io/gorm"

type Consignee struct {
	Email string
	StreetAddress string
	City	string
	State	string
	Country	string
	ZipCode int32
}

type Order struct {
	gorm.Model
	OrderId string `gorm:"type:varchar(100);uniqueIndex"`
	UserId	uint32	`gorm:"type:int(11)"`
	UserCurrency string	`gorm:"type:varchar(10)"`
	Consignee Consignee	`gorm:"embedded"`
	OrderItem []OrderItem `gorm:"foreignKey:OrderIdRefer;references:OrderId"`
}

func (Order) TableName() string {
	return "order"
}

type OrderQuery struct {
	db *gorm.DB
}

func NewOrderQuery(db *gorm.DB) *OrderQuery {
	return &OrderQuery{
		db: db,
	}
}

func (q *OrderQuery) ListOrder(userId uint32) ([]*Order, error) {
	var orders []*Order
	err := q.db.Where("user_id = ?", userId).Preload("OrderItem").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

