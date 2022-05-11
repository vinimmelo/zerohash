package entities

import "math/big"

type CurrencyQueue struct {
	elements  []*big.Float
	maxNumber int // Defines the maximum number of the queue
}

func NewCurrencyQueue(maxNumber int) *CurrencyQueue {
	return &CurrencyQueue{
		maxNumber: maxNumber,
	}
}

func (q *CurrencyQueue) IsEmpty() bool {
	return len(q.elements) == 0
}

// Enqueue an element, if the queue reaches the maximum
// length defined, it will drop (dequeue) the first element of it
func (q *CurrencyQueue) Enqueue(currency *big.Float) {
	if q.ReachMaximum() {
		q.Dequeue()
	}
	q.elements = append(q.elements, currency)
}

// Drop (dequeue) the first element of the queue
func (q *CurrencyQueue) Dequeue() *big.Float {
	var first *big.Float
	if !q.IsEmpty() {
		first = q.elements[0]
		q.elements = q.elements[1:]
	}
	return first
}

// Verify if the queue reached the maximum size defined
func (q *CurrencyQueue) ReachMaximum() bool {
	return len(q.elements) == q.maxNumber
}

func (q *CurrencyQueue) Elements() []*big.Float {
	return q.elements
}

// Sums all the elements and returns the value
func (q *CurrencyQueue) SumElements() (sum *big.Float) {
	sum = big.NewFloat(0)
	for _, element := range q.Elements() {
		sum.Add(sum, element)
	}
	return
}
