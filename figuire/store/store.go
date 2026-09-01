package store

import (
	"errors"
	"figuire/figuire/merch"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Store struct {
	db *gorm.DB
}

func NewStore(path string) (*Store, error) {
	dsn := path + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Product{}, &Watch{}, &PriceEvent{}); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Record(item merch.Item) (*Change, error) {
	var p Product
	err := s.db.Where("site = ? AND code = ?", item.Site, item.Code).First(&p).Error
	isNew := errors.Is(err, gorm.ErrRecordNotFound)

	if err != nil && !isNew {
		return nil, err
	}

	c := &Change{
		IsNew:           isNew,
		NewPriceJPY:     item.PriceJPY,
		NewAvailability: item.Availability,
	}

	if !isNew {
		c.OldPriceJPY = p.PriceJPY
		c.OldAvailability = merch.Availability(p.Availability)
	}

	p.Site, p.Code, p.Name = item.Site, item.Code, item.Name
	p.Maker, p.JanCode, p.URL = item.Maker, item.JanCode, item.URL
	p.ImageURL, p.ReleaseText, p.PreOwned = item.ImageURL, item.ReleaseText, item.PreOwned
	p.PriceJPY, p.ListPriceJPY, p.BuyPriceJPY = item.PriceJPY, item.ListPriceJPY, item.BuyPriceJPY
	p.Availability = string(item.Availability)
	p.LastCheckedAt = item.CheckedAt

	if isNew {
		p.FirstSeenAt = item.CheckedAt
	}

	if c.Any() {
		p.LastChangedAt = item.CheckedAt
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&p).Error; err != nil {
			return err
		}

		if c.Any() {
			return tx.Create(&PriceEvent{
				ProductID:    p.ID,
				ObservedAt:   item.CheckedAt,
				PriceJPY:     item.PriceJPY,
				Availability: string(item.Availability),
			}).Error
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	c.Product = &p
	return c, nil
}

func (s *Store) AddWatch(w *Watch) error {
	err := s.db.Create(w).Error
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return ErrAlreadyWatching
	}
	return err
}

func (s *Store) RemoveWatch(channelID string, productID uint) (bool, error) {
	res := s.db.Where("channel_id = ? AND product_id = ?", channelID, productID).Delete(&Watch{})
	return res.RowsAffected > 0, res.Error
}

func (s *Store) WatchedProducts() ([]Product, error) {
	var products []Product
	err := s.db.Joins("JOIN watches ON watches.product_id = products.id").
		Group("products.id").
		Order("products.name asc").
		Find(&products).Error
	return products, err
}

func (s *Store) WatchesFor(productID uint) ([]Watch, error) {
	var watches []Watch
	err := s.db.Where("product_id = ?", productID).Find(&watches).Error
	return watches, err
}

func (s *Store) WatchesInChannel(channelID string) ([]Watch, error) {
	var watches []Watch
	err := s.db.Where("channel_id = ?", channelID).Find(&watches).Error
	return watches, err
}

func (s *Store) ProductByCode(site, code string) (*Product, error) {
	var p Product
	if err := s.db.Where("site = ? AND code = ?", site, code).First(&p).Error; err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *Store) History(productID uint) ([]PriceEvent, error) {
	var events []PriceEvent
	err := s.db.Where("product_id = ?", productID).Order("observed_at asc").Find(&events).Error
	return events, err
}
