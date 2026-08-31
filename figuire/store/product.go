package store

import "time"

type Product struct {
	ID   uint   `gorm:"primaryKey"`
	Site string `gorm:"uniqueIndex:idx_site_code;not null"`
	Code string `gorm:"uniqueIndex:idx_site_code;not null"`

	Name        string
	Maker       string
	JanCode     string `gorm:"index"`
	URL         string
	ImageURL    string
	ReleaseText string
	PreOwned    bool

	PriceJPY     int
	ListPriceJPY int
	BuyPriceJPY  int

	Availability string

	FirstSeenAt   time.Time
	LastCheckedAt time.Time
	LastChangedAt time.Time
}
