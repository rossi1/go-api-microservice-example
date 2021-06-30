package handler

import (
	"time"
)

type CategorySerializer struct {
	Name    string              `json:"name"`
	Product []ProductSerializer `json:"product"`
}

type ProductSerializer struct {
	Name        string     `json:"name"`
	Tax         *string    `json:"tax"`
	Description *string    `json:"description"`
	Weight      *string    `json:"weight"`
	Expires     *time.Time `json:"expires"`
	BarCode     *string    `json:"bar_code"`
	Discount    *string    `json:"discount"`
	Image       *string    `json:"image"`
}
