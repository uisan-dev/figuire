package amiami

import (
	"figuire/figuire/merch"
	"time"
)

type apiResponse struct {
	RSuccess bool     `json:"RSuccess"`
	RMessage string   `json:"RMessage"`
	Item     *apiItem `json:"item"`
}

type apiItem struct {
	GCode       string `json:"gcode"`
	SCode       string `json:"scode"`
	SNameSimple string `json:"sname_simple"`
	JanCode     string `json:"jancode"`
	MakerName   string `json:"maker_name"`
	ThumbURL    string `json:"thumb_url"`
	ReleaseDate string `json:"releasedate"`

	Price     int `json:"price"`
	ListPrice int `json:"list_price"`
	BuyPrice  int `json:"buy_price"`

	CartType      int `json:"cart_type"`
	ConditionFlag int `json:"condition_flag"`
}

func (ai *apiItem) availability() merch.Availability {
	switch ai.CartType {
	case 9:
		return merch.InStock
	case 8:
		return merch.Preorder
	case 7:
		return merch.Backorder
	case 0, 1, 2, 3, 4, 5, 6:
		return merch.Unavailable
	default:
		return merch.Unknown
	}
}

func (ai *apiItem) toMerchItem() merch.Item {
	return merch.Item{
		Site:         "amiami",
		Code:         ai.GCode,
		Name:         ai.SNameSimple,
		URL:          "https://www.amiami.com/eng/detail/?gcode=" + ai.GCode,
		ImageURL:     ImageBase + ai.ThumbURL,
		JanCode:      ai.JanCode,
		Maker:        ai.MakerName,
		PriceJPY:     ai.Price,
		ListPriceJPY: ai.ListPrice,
		BuyPriceJPY:  ai.BuyPrice,
		PreOwned:     ai.ConditionFlag == 1,
		Availability: ai.availability(),
		ReleaseText:  ai.ReleaseDate,
		CheckedAt:    time.Now(),
	}
}
