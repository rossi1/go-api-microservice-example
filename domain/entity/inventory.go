package entity

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID uuid.UUID `gorm:"id;type:uuid;primary_key" json:"id,omitempty"`
}

type Category struct {
	Base
	Name      string         `gorm:"size:200;name;unique" json:"name" validate:"email" message:"email is invalid"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP;created_at" json:"created_at"`
	Deleted   gorm.DeletedAt `json:"-"`
	Products  []Product      `gorm:"constraint:foreignKey:CategoryID;OnDelete:CASCADE;products" json:"products"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {

	id, err := uuid.NewV4()

	if err != nil {
		return err
	}

	c.ID = id
	return nil
}

type Product struct {
	Base
	Name        string         `gorm:"size:200;name;unique" json:"name" validate:"email" message:"email is invalid"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP;created_at" json:"created_at"`
	Deleted     gorm.DeletedAt `json:"-"`
	Tax         *string        `gorm:"size:200;null;name" json:"tax"`
	Description *string        `gorm:"size:200;null;name" json:"description"`
	Weight      *string        `gorm:"size:200;null;name" json:"weight"`
	Expires     *time.Time     `gorm:"size:200;null;name" json:"expires"`
	BarCode     *string        `gorm:"size:200;null;name" json:"bar_code"`
	Discount    *string        `gorm:"size:200;null;name" json:"discount"`
	Image       *string        `gorm:"size:200;null;name" json:"image"`
	CategoryID  uuid.UUID      `gorm:"CategoryID" json:"-"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {

	id, err := uuid.NewV4()

	if err != nil {
		return err
	}

	p.ID = id
	return nil
}
