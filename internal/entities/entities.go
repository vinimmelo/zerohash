package entities

import (
	"math/big"
)

type TickerSymbol string

const (
	BTCUSD TickerSymbol = "BTC-USD"
	ETHUSD TickerSymbol = "ETH-USD"
	ETHBTC TickerSymbol = "ETH-BTC"
)

type TickerInfo struct {
	TickerSymbol TickerSymbol
	Price        *big.Float
}
