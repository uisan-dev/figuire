package merch

import "time"

type Availability string

const (
	InStock     Availability = "in_stock"
	Preorder    Availability = "preorder"
	Backorder   Availability = "backorder"
	Unavailable Availability = "unavailable"
	Unknown     Availability = "unknown"
)

var displayStrings = map[Availability]string{
	InStock:     "In stock",
	Preorder:    "Preorder",
	Backorder:   "Backorder",
	Unavailable: "Unavailable",
	Unknown:     "Unknown",
}

func (a Availability) AsDisplayString() string {
	return displayStrings[a]
}

type Item struct {
	Site         string
	Code         string
	Name         string
	URL          string
	ImageURL     string
	JanCode      string
	Maker        string
	PriceJPY     int
	ListPriceJPY int
	BuyPriceJPY  int
	PreOwned     bool
	Availability Availability
	ReleaseText  string
	CheckedAt    time.Time
}
