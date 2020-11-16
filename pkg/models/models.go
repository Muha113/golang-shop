package models

import (
	"time"
)

type Product struct {
	ID          int
	ShopID      int
	Name        string
	Description string
	Price       float32
}

type Shop struct {
	ID       int
	Name     string
	URL      string
	Open     time.Time
	Close    time.Time
	Products []*Product
}
