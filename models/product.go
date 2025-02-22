package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string
	Description string
	Picture     string
	Price       float32

	Categories []Category `json:"categories" gorm:"many2many:product_category"`
}

func (p Product) TableName() string {
	return "product"
}

type ProductQuery struct {
	db *gorm.DB
}

func (q *ProductQuery) GetById(productId int32) (product Product, err error) {
	err = q.db.Model(&Product{}).First(&product, productId).Error
	return
}

func (q *ProductQuery) ProductExisted(productName string) (exist bool, err error) {
	var cnt int64
	err = q.db.Model(&Product{}).Where("name=?", productName).Count(&cnt).Error
	exist = cnt > 0
	return
}

func (q *ProductQuery) SearchProducts(query string) (products []*Product, err error) {
	err = q.db.Model(&Product{}).Find(&products, "name like ? or description like ?",
		"%"+query+"%", "%"+query+"%",
	).Error
	return
}

func NewProductQuery(db *gorm.DB) *ProductQuery {
	return &ProductQuery{
		db: db,
	}
}
