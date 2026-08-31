package figurebot

import "time"

type Product struct {
	ID   uint   `gorm:"primaryKey"`
	Site string `gorm:"uniqueIndex:idx_site_code;not null"`
	Code string `gorm:"uniqueIndex:idx_site_code;not null"`
}

type Watch struct {
	ID        uint   `gorm:"primaryKey"`
	GuildID   string `gorm:"index;not null"`
	ChannelID string `gorm:"uniqueIndex:idx_channel_product;not null"`
	ProductID string `gorm:"uniqueIndex:idx_channel_product;not null"`
	AddedBy   string
	CreatedAt time.Time
}
