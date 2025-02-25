package models

import (
	"errors"
	"goshop/pkg/redis"

	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserId     uint64 `gorm:"type:int(11);not null;index:idx_user_id"` // Create index here.
	ProductId  uint64 `gorm:"type:int(11);not null;"`
	ProductCnt uint64 `gorm:"type:int(11);not null;"`
}

func (c Cart) TableName() string {
	return "cart"
}

type CartQuery struct {
	db  *gorm.DB
	rdb *redis.Rdb
}

func NewCartQuery(db *gorm.DB) *CartQuery {
	return &CartQuery{
		db: db,
	}
}

func (q *CartQuery) GetByUserId(user_id uint32) (carts []*Cart, err error) {
	err = q.db.Model(&Cart{}).Where(&Cart{UserId: uint64(user_id)}).Find(&carts).Error
	return
}

func (q *CartQuery) AddProduct(item *Cart) error {
	var row Cart
	err := q.db.Model(&Cart{}).
		Where(&Cart{UserId: item.UserId, ProductId: item.ProductId}).
		First(&row).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Add the item if record not found.
		return err
	}
	if row.ID > 0 {
		return q.db.Model(&Cart{}).
			Where(&Cart{UserId: item.UserId, ProductId: item.ProductId}).
			UpdateColumn("product_cnt", gorm.Expr("product_cnt+?", item.ProductCnt)).Error
	}
	return q.db.Create(item).Error
}

func (q *CartQuery) CleanByUserId(user_id uint32) error {
	if user_id == 0 {
		return errors.New("user id is required.")
	}
	return q.db.Delete(&Cart{}, "user_id = ?", user_id).Error
}
