package store

import "figuire/figuire/merch"

type Change struct {
	Product *Product
	IsNew   bool

	OldPriceJPY     int
	NewPriceJPY     int
	OldAvailability merch.Availability
	NewAvailability merch.Availability
}

func (c *Change) PriceChanged() bool {
	return !c.IsNew && c.OldPriceJPY != c.NewPriceJPY
}

func (c *Change) AvailabilityChanged() bool {
	return !c.IsNew && c.OldAvailability != c.NewAvailability
}

func (c *Change) BecameAvailable() bool {
	return !c.IsNew && c.OldAvailability != merch.InStock && c.NewAvailability == merch.InStock
}

func (c *Change) Any() bool {
	return c.IsNew || c.AvailabilityChanged() || c.PriceChanged()
}
