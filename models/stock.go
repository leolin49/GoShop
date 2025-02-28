package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Stock struct {
	gorm.Model
	ProductId uint32
	Count     uint64
}

func (s Stock) TableName() string {
	return "stock"
}

type StockQuery struct {
	db *gorm.DB
}

func NewStockQuery(db *gorm.DB) *StockQuery {
	return &StockQuery{
		db: db,
	}
}

func (q *StockQuery) GetStock(productId uint32) (count uint64, err error) {
	err = q.db.Model(&Stock{}).Select("count").Where("product_id = ?", productId).First(&count).Error
	return
}

func (q *StockQuery) AddStock(productId uint32, addCount uint64) (err error) {
	// NOTE: Don't use the return value 'count' here cause Concurrency safe.
	_, err = q.GetStock(productId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Add count if record not fount
		return q.db.Create(&Stock{
			ProductId: productId,
			Count:     addCount,
		}).Error
	}
	// NOTE: Must use 'Update...' to add line lock in the table.
	return q.db.Model(&Stock{}).
		Where("product_id = ?", productId).
		UpdateColumn("count", gorm.Expr("count + ?", addCount)).Error
}

// NOTE: Use the atomicity of transactions and optimistic locks to
// avoid oversold problems.
// 1. If count>=subCount, Update is added
// 2. If no operation is performed by Update, the stock is insufficient
func (q *StockQuery) SubStock(productId uint32, subCount uint64) (err error) {
	_, err = q.GetStock(productId)
	if err != nil {
		return
	}
	return q.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Stock{}).
			Where("product_id = ? and count >= ?", productId, subCount).
			UpdateColumn("count", gorm.Expr("count - ?", subCount))
		if res.RowsAffected == 0 {
			// no any row affected becasue count < subCount.
			return fmt.Errorf("no more stock for product [%d]", productId)
		}
		// success
		return nil
	})
}
