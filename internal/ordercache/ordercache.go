package ordercache

import (
	"fmt"
	"log"
	"sync"

	"wblayerzero/internal/entities"
)

type (
	OrderCache struct {
		m *sync.Map

		logger *log.Logger
	}

	Key string
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
		oc.m.Store(Key(order.UID), &order)
	}

	return nil
}

// GetOne returns the order with a given orderUID
// Method doesn't include any checks of type upcasting, because the Upload completes all the checks needed while dumping elems into cache
func (oc *OrderCache) GetSingle(orderUID string) (*entities.Order, error) {
	preOrder, ok := oc.m.Load(Key(orderUID))
	if !ok {
		return nil, fmt.Errorf("order not found for the key %q", orderUID)
	}

	return preOrder.(*entities.Order), nil
}

// GetAll returns all the orders included within oc
// You mustn't count on consistency of the internal m
func (oc *OrderCache) GetAll() []*entities.Order {
	const (
		defaultStartSize = 1 << 6
	)

	var (
		orders []*entities.Order
	)

	orders = make([]*entities.Order, defaultStartSize)

	oc.m.Range(func(key, value any) bool {
		orders = append(orders, value.(*entities.Order))
		return true
	})

	return orders
}
