package store

import "time"

type Watch struct {
	ID        uint    `gorm:"primaryKey"`
	GuildID   string  `gorm:"index;not null"`
	ChannelID string  `gorm:"uniqueIndex:idx_channel_product;not null"`
	ProductID uint    `gorm:"uniqueIndex:idx_channel_product;not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
	AddedBy   string
	CreatedAt time.Time
}
