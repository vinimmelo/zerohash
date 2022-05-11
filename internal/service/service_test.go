package service

import (
	"io"
	"math/big"
	"testing"

	"bytes"

	"fmt"

	"github.com/vinimmelo/zerohash/internal/entities"
	"github.com/vinimmelo/zerohash/internal/logger"
)

func Test_ExecuteVWAP(t *testing.T) {
	log := logger.New(logger.DEBUG)
	service := New(log)

	t.Run("when the ticker is BTC-USD", func(t *testing.T) {
		tickerInfo := &entities.TickerInfo{
			TickerSymbol: entities.BTCUSD,
			Price:        big.NewFloat(31565.54),
		}

		result := service.ExecuteVWAP(tickerInfo, nil)

		if result != 31565.54 {
			t.Errorf("should be 31565.54, but it is: %f", result)
		}

		// after three more executions
		tickerInfoTwo := &entities.TickerInfo{
			TickerSymbol: entities.BTCUSD,
			Price:        big.NewFloat(31565.49),
		}
		tickerInfoThree := &entities.TickerInfo{
			TickerSymbol: entities.BTCUSD,
			Price:        big.NewFloat(31565.39),
		}
		tickerInfoFour := &entities.TickerInfo{
			TickerSymbol: entities.BTCUSD,
			Price:        big.NewFloat(31565.47),
		}
		service.ExecuteVWAP(tickerInfoTwo, nil)
		service.ExecuteVWAP(tickerInfoThree, nil)
		result = service.ExecuteVWAP(tickerInfoFour, nil)

		if result != 31565.4725 {
			t.Errorf("should be 31565.4725, but it is: %f", result)
		}

		t.Run("when more than 1000 tickets are executed", func(t *testing.T) {
			var result float64
			for i := 1; i <= 1000; i++ {
				var ticker *entities.TickerInfo
				if i <= 800 {
					ticker = &entities.TickerInfo{
						TickerSymbol: entities.BTCUSD,
						Price:        big.NewFloat(999.99),
					}
				} else if i > 800 && i <= 900 {
					ticker = &entities.TickerInfo{
						TickerSymbol: entities.BTCUSD,
						Price:        big.NewFloat(3999.37),
					}
				} else if i > 900 && i <= 1000 {
					ticker = &entities.TickerInfo{
						TickerSymbol: entities.BTCUSD,
						Price:        big.NewFloat(3999.79),
					}
				}
				result = service.ExecuteVWAP(ticker, nil)
			}

			if result != 3999.58 {
				t.Errorf("should be 3999.58, not %f", result)
			}
		})

		t.Run("when it prints the right thing", func(t *testing.T) {
			// Reset the queue
			Queues[entities.BTCUSD] = entities.NewCurrencyQueue(200)

			tests := []struct {
				tickerInfo    *entities.TickerInfo
				expectedPrint string
			}{
				{
					tickerInfo: &entities.TickerInfo{
						TickerSymbol: entities.BTCUSD,
						Price:        big.NewFloat(31565.54),
					},
					expectedPrint: fmt.Sprintf("Ticker: BTC-USD\tVWAP: %.6f\n", 31565.54),
				},
				{
					tickerInfo: &entities.TickerInfo{
						TickerSymbol: entities.BTCUSD,
						Price:        big.NewFloat(31565.49),
					},
					expectedPrint: fmt.Sprintf("Ticker: BTC-USD\tVWAP: %.6f\n", 31565.515),
				},
				{
					tickerInfo: &entities.TickerInfo{
						TickerSymbol: entities.BTCUSD,
						Price:        big.NewFloat(31565.37),
					},
					expectedPrint: fmt.Sprintf("Ticker: BTC-USD\tVWAP: %.6f\n", 31565.466666666666),
				},
			}

			var b bytes.Buffer
			for _, test := range tests {
				service.ExecuteVWAP(test.tickerInfo, []io.Writer{&b})

				result := string(b.String())

				if result != test.expectedPrint {
					t.Errorf("printed the wrong thing, result: <%s> expected: <%s>", result, test.expectedPrint)
				}
				b.Reset()
			}
		})
	})

	t.Run("when the ticker is ETH-USD", func(t *testing.T) {
		var b bytes.Buffer

		tickerInfo := &entities.TickerInfo{
			TickerSymbol: entities.ETHUSD,
			Price:        big.NewFloat(31565.54),
		}

		service.ExecuteVWAP(tickerInfo, []io.Writer{&b})
		result := string(b.String())
		expected := fmt.Sprintf("Ticker: ETH-USD\tVWAP: %.6f\n", 31565.54)

		if result != expected {
			t.Errorf("printed the wrong thing, result: <%s> expected: <%s>", result, expected)
		}

		tickerInfoTwo := &entities.TickerInfo{
			TickerSymbol: entities.ETHUSD,
			Price:        big.NewFloat(31565.49),
		}

		b.Reset()
		service.ExecuteVWAP(tickerInfoTwo, []io.Writer{&b})
		result = string(b.String())
		expected = fmt.Sprintf("Ticker: ETH-USD\tVWAP: %.6f\n", 31565.515)

		if result != expected {
			t.Errorf("printed the wrong thing, result: <%s> expected: <%s>", result, expected)
		}
	})

	t.Run("when the ticker is ETH-BTC", func(t *testing.T) {
		var b bytes.Buffer

		tickerInfo := &entities.TickerInfo{
			TickerSymbol: entities.ETHBTC,
			Price:        big.NewFloat(31565.54),
		}

		service.ExecuteVWAP(tickerInfo, []io.Writer{&b})
		result := string(b.String())
		expected := fmt.Sprintf("Ticker: ETH-BTC\tVWAP: %.6f\n", 31565.54)

		if result != expected {
			t.Errorf("printed the wrong thing, result: <%s> expected: <%s>", result, expected)
		}

		tickerInfoTwo := &entities.TickerInfo{
			TickerSymbol: entities.ETHBTC,
			Price:        big.NewFloat(31565.49),
		}

		b.Reset()
		service.ExecuteVWAP(tickerInfoTwo, []io.Writer{&b})
		result = string(b.String())
		expected = fmt.Sprintf("Ticker: ETH-BTC\tVWAP: %.6f\n", 31565.515)

		if result != expected {
			t.Errorf("printed the wrong thing, result: <%s> expected: <%s>", result, expected)
		}
	})
}
