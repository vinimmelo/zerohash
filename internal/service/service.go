package service

import (
	"fmt"
	"io"
	"math/big"

	"math"

	"github.com/vinimmelo/zerohash/internal/entities"
	"github.com/vinimmelo/zerohash/internal/logger"
)

var Queues map[entities.TickerSymbol]*entities.CurrencyQueue

func init() {
	// Create short-live queues
	Queues = make(map[entities.TickerSymbol]*entities.CurrencyQueue)
	Queues[entities.BTCUSD] = entities.NewCurrencyQueue(200)
	Queues[entities.ETHUSD] = entities.NewCurrencyQueue(200)
	Queues[entities.ETHBTC] = entities.NewCurrencyQueue(200)
}

type Service struct {
	log logger.Logger
}

func New(log logger.Logger) *Service {
	return &Service{log: log}
}

// Main function, it will execute the VWAP calculation engine and print the result
// to the configured printers/writers
func (s *Service) ExecuteVWAP(tickerInfo *entities.TickerInfo, printers []io.Writer) (result float64) {
	queue := Queues[tickerInfo.TickerSymbol]
	queue.Enqueue(tickerInfo.Price)

	result = s.roundedResult(queue)

	for _, writer := range printers {
		s.printVWAP(result, tickerInfo.TickerSymbol, writer)
	}
	return
}

// This function should sum all the elements and divides them by their queue size.
// Also, this function rounds the result to the sixth decimal place
func (s *Service) roundedResult(queue *entities.CurrencyQueue) float64 {
	sum := queue.SumElements()

	// Using big.Float improves precision, but it's still inaccurate.
	// The right way to do it is to create or import a library that
	// deals correctly with decimals
	result := sum.Quo(sum, big.NewFloat(float64(len(queue.Elements()))))
	floatResult, _ := result.Float64()
	ratio := math.Pow(10, float64(6))
	return math.Round(floatResult*ratio) / ratio
}

func (s *Service) printVWAP(value float64, tickerSymbol entities.TickerSymbol, w io.Writer) {
	fmt.Fprintf(w, "Ticker: %s\tVWAP: %.6f\n", tickerSymbol, value)
}
