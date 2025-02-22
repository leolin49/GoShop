package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string
	Description string

	Products []Product `json:"product" gorm:"many2many:product_category"`
}

func (c Category) TableName() string {
	return "category"
}

type CategoryQuery struct {
	db *gorm.DB
}

func (q *CategoryQuery) GetProductsByCategoryName(name string) (categories []Category, err error) {
	err = q.db.Model(&Category{}).Where(&Category{Name: name}).Find(&categories).Error
	return
}

func (q *CategoryQuery) GetIdsByNames(names []string) (categoryIds []int64, err error) {
	for _, name := range names {
		id, err := q.GetIdByName(name)
		if err != nil {
			break
		}
		categoryIds = append(categoryIds, id)
	}
	return
}

// SELECT id FROM category WHERE name=`name`;
func (q *CategoryQuery) GetIdByName(name string) (categoryId int64, err error) {
	err = q.db.Model(&Category{}).Select("id").Where("name=?", name).Find(&categoryId).Error
	return
}

func (q *CategoryQuery) ListProducts(categoryName string, page int, pageSize int) (products []*Product, err error) {
	var category Category
	err = q.db.Where("name=?", categoryName).Preload("Products", func(db *gorm.DB) *gorm.DB {
		return db.Model(&Product{}).Limit(pageSize).Offset((page - 1) * pageSize)
	}).Find(&category).Error
	for _, pd := range category.Products {
		products = append(products, &pd)
	}
	return
}

func NewCategoryQuery(db *gorm.DB) *CategoryQuery {
	return &CategoryQuery{
		db: db,
	}
}
