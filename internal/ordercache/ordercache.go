package ordercache

import (
	"fmt"
	"log"
	"sync"

	"wblayerzero/internal/entities"
)

type (
	Key string

	OrderCache struct {
		m      *sync.Map
		logger *log.Logger
	}
)

// New returns a new initialized OrderCache instance with the given logger
func New(logger *log.Logger) *OrderCache {
	return &OrderCache{m: &sync.Map{}, logger: logger}
}

// Upload stores the given orders in internal field m of oc
func (oc *OrderCache) Upload(orders ...entities.Order) error {
	if len(orders) == 0 {
		return fmt.Errorf("no values given")
	}

	for _, order := range orders {
		oc.m.Store(Key(order.UID), order)
	}

	return nil
}

// Get returns the order with a given orderUID
func (oc *OrderCache) Get(orderUID string) (entities.Order, error) {
	preOrder, ok := oc.m.Load(Key(orderUID))
	if !ok {
		return entities.Order{}, fmt.Errorf("order not found for the key %q", orderUID)
	}

	v := preOrder.(entities.Order)
	return v, nil
}

// DownloadOrders prints all the orders from the internal m
func (oc *OrderCache) DownloadOrders() {
	oc.m.Range(func(key, value any) bool {
		fmt.Println(key.(Key), value.(entities.Order))
		return true
	})
}
