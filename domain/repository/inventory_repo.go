package repository

import (
	"errors"
	"fmt"
	"github/rossi1/go-api-microservice-example/domain/entity"
	"github/rossi1/go-api-microservice-example/internal/pub"
	"github/rossi1/go-api-microservice-example/internal/search"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 5
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

var (
	ErrNotFound = errors.New("record not found")
	ErrDelete   = errors.New("")
	ErrUpdate   = errors.New("unable to update record")
	ErrFetch    = errors.New("unable to fetch records")
	ErrCreate   = errors.New("unable to create record")
)

type CategoryRepository interface {
	Create(name string, product []entity.Product) (*entity.Category, error)
	Delete(id uuid.UUID) error
	Find(id uuid.UUID) (*entity.Category, error)
	FindAll(r *http.Request) ([]entity.Category, error)
	Update(id uuid.UUID, name string) (*entity.Category, error)
}

type ProductRepository interface {
	Create(product entity.Product, categoryID uuid.UUID) (*entity.Product, error)
	Delete(id uuid.UUID) error
	Find(id uuid.UUID) (*entity.Product, error)
	Update(id uuid.UUID, params map[string]interface{}) (*entity.Product, error)
	FindAll(id uuid.UUID, r *http.Request) ([]entity.Product, error)
}

type SearchRepository interface {
	search.Searcher
}

type PubRepository interface {
	pub.Pub
}

type categoryRepo struct {
	db *gorm.DB
	//search   SearchRepository
	//producer PubRepository
}

func (c *categoryRepo) Create(name string, product []entity.Product) (*entity.Category, error) {
	category := entity.Category{Name: name}
	if len(product) > 0 {
		category.Products = product
	}
	err := c.db.Create(&category).Error
	if err != nil {
		log.Println(err.Error())
		return nil, ErrCreate
	}
	return &category, nil

}

func (c *categoryRepo) Delete(id uuid.UUID) error {
	var category entity.Category

	err := c.db.Where("id = ?", id).Delete(&category).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println(err.Error())
		return ErrDelete
	}
	return nil
}

func (c *categoryRepo) Find(id uuid.UUID) (*entity.Category, error) {
	var category entity.Category
	err := c.db.Preload("Products").First(&category, "id = ?", id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println(err.Error())
		return nil, ErrNotFound
	}
	fmt.Println(category)
	return &category, nil
}

func (c *categoryRepo) Update(id uuid.UUID, name string) (*entity.Category, error) {
	var category entity.Category
	err := c.db.First(&category, "id = ?", id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println(err.Error())
		return nil, ErrNotFound
	}

	err = c.db.Model(&category).Where("id = ?", id).Update("name", name).Error
	if err != nil {
		log.Println(err.Error())
		return nil, ErrUpdate
	}

	return &category, nil
}

func (c *categoryRepo) FindAll(r *http.Request) ([]entity.Category, error) {
	var categories []entity.Category
	err := c.db.Scopes(Paginate(r)).Preload("Products").Find(&categories).Error
	if err != nil {
		log.Println(err.Error())
		return nil, ErrFetch
	}
	return categories, nil

}

func NewCategoryRepo(db *gorm.DB) CategoryRepository {
	return &categoryRepo{db: db}
}

type productRepo struct {
	db *gorm.DB
	//search   SearchRepository
	//producer PubRepository
}

func (p *productRepo) Create(product entity.Product, categoryID uuid.UUID) (*entity.Product, error) {
	var category entity.Category
	err := p.db.First(&category, "id = ?", categoryID).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	category.Products = []entity.Product{product}
	if err = tx.Save(&category).Error; err != nil {
		return nil, err
	}

	product.CategoryID = category.ID
	if err = tx.Create(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *productRepo) Delete(id uuid.UUID) error {
	var product entity.Product
	err := p.db.Where("id = ?", id).Delete(&product).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println(err.Error())
		return ErrDelete
	}
	return nil
}

func (p *productRepo) Find(id uuid.UUID) (*entity.Product, error) {
	var product entity.Product
	err := p.db.First(&product, "id = ?", id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println(err.Error())
		return nil, ErrFetch
	}
	return &product, nil

}

func (p *productRepo) Update(id uuid.UUID, params map[string]interface{}) (*entity.Product, error) {
	var product entity.Product

	err := p.db.Model(&product).Where("id = ?", id).Error
	if err != nil {
		log.Println(err.Error())
		return nil, ErrUpdate
	}
	if name := params["name"].(string); name != "" {
		product.Name = name
	}

	if tax := params["tax"].(*string); tax != nil {
		product.Tax = tax
	}

	if description := params["description"].(*string); description != nil {
		product.Description = description
	}

	if weight := params["weight"].(*string); weight != nil {
		product.Weight = weight
	}

	if expires := params["expires"].(*time.Time); expires != nil {
		product.Expires = expires
	}

	if barcode := params["barcode"].(*string); barcode != nil {
		product.BarCode = barcode
	}

	if discount := params["discount"].(*string); discount != nil {
		product.Discount = discount
	}

	if image := params["image"].(*string); image != nil {
		product.Image = image
	}
	p.db.Save(&product)
	return &product, nil
}

func (c *productRepo) FindAll(id uuid.UUID, r *http.Request) ([]entity.Product, error) {
	var products []entity.Product
	err := c.db.Scopes(Paginate(r)).Where("category_id = ?", id).Find(&products).Error
	if err != nil {
		fmt.Println(err.Error())
		return nil, ErrFetch
	}
	return products, nil

}

func NewProductRepo(db *gorm.DB) ProductRepository {
	return &productRepo{db: db}
}
