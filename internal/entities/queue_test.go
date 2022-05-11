package entities

import (
	"math/big"
	"testing"
)

func Test_CurrencyQueue_IsEmpty(t *testing.T) {
	queue := NewCurrencyQueue(200)

	t.Run("when the queue is recently created", func(t *testing.T) {
		if !queue.IsEmpty() {
			t.Errorf("queue should be empty")
		}
	})

	t.Run("when the queue dequeue all its elements", func(t *testing.T) {
		currency := big.NewFloat(100.65)
		queue.Enqueue(currency)

		if queue.IsEmpty() {
			t.Errorf("queue shouldn't be empty")
		}

		queue.Dequeue()
		if !queue.IsEmpty() {
			t.Error("queue should be empty")
		}
	})
}

func Test_CurrencyQueue_Enqueue(t *testing.T) {
	queue := NewCurrencyQueue(200)
	currency := big.NewFloat(100.0)

	queue.Enqueue(currency)

	if len(queue.Elements()) != 1 {
		t.Errorf("queue should have 1 element")
	}

	currencyTwo := big.NewFloat(100.0)
	queue.Enqueue(currencyTwo)

	if len(queue.Elements()) != 2 {
		t.Errorf("queue should have 2 elements")
	}

	t.Run("when the queue reaches the maximum, it removes the first element automatically", func(t *testing.T) {
		queue = NewCurrencyQueue(200)

		for i := 0; i < 200; i++ {
			currency := big.NewFloat(float64(i))
			queue.Enqueue(currency)
			if i == 199 && !queue.ReachMaximum() {
				t.Error("should've reached the maximum after 200 elements enqueued")
			} else if i < 199 && queue.ReachMaximum() {
				t.Errorf("shouldn't reach the maximum before 200 elements, index: %d", i)
			}
		}

		queue.Enqueue(big.NewFloat(100.0))
		queue.Enqueue(big.NewFloat(100.0))
		queue.Enqueue(big.NewFloat(100.0))

		if len(queue.Elements()) != 200 {
			t.Errorf("queue should have 200 elements, but has %d", len(queue.Elements()))
		}
	})
}

func Test_CurrencyQueue_Dequeue(t *testing.T) {
	queue := NewCurrencyQueue(200)

	t.Run("when the queue is empty, should return nil", func(t *testing.T) {
		result := queue.Dequeue()

		if result != nil {
			t.Errorf("should return an empty Currency, returned: %s", result.String())
		}
	})

	t.Run("when the queue has elements, should return the first one", func(t *testing.T) {
		currency := big.NewFloat(100.0)
		currencyTwo := big.NewFloat(400.0)
		currencyThree := big.NewFloat(300.0)

		queue.Enqueue(currency)
		queue.Enqueue(currencyTwo)
		queue.Enqueue(currencyThree)

		result := queue.Dequeue()
		if result.Cmp(big.NewFloat(100.0)) != 0 {
			t.Errorf("dequeue should return the first element with value of 100.0, returned: %s", result)
		}
	})
}

func Test_CurrencyQueue_SumElements(t *testing.T) {
	queue := NewCurrencyQueue(200)

	for i := 1; i < 100; i++ {
		queue.Enqueue(big.NewFloat(float64(i)))
	}
	expected := big.NewFloat(4950)

	if sum := queue.SumElements(); sum.Cmp(expected) != 0 {
		t.Errorf("the sum of the elements should've be %f, and the result was %f", expected, sum)
	}
}
