package websocket

import (
	"math/big"
	"strconv"

	"github.com/vinimmelo/zerohash/internal/entities"
)

type SubscribeResponse struct {
	Type      string `json:"type"`
	ProductID string `json:"product_id"`
	Price     string `json:"price"`
}

func (s SubscribeResponse) ToTickerInfo() (*entities.TickerInfo, error) {
	ticker := entities.TickerSymbol(s.ProductID)
	value, err := strconv.ParseFloat(s.Price, 64)
	if err != nil {
		return nil, err
	}

	price := big.NewFloat(value)

	return &entities.TickerInfo{
		TickerSymbol: ticker,
		Price:        price,
	}, nil
}
