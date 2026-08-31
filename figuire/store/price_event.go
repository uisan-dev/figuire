package store

import "time"

type PriceEvent struct {
	ID           uint      `gorm:"primaryKey"`
	ProductID    uint      `gorm:"index:idx_product_time,priority:1;not null"`
	ObservedAt   time.Time `gorm:"index:idx_product_time,priority:2;not null"`
	PriceJPY     int
	Availability string
}
