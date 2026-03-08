package order

import (
	"errors"

	"github.com/DiegoAmin/AmazonClone_PAP/internal/product"
)

// OrderItem represents an item in the order, which includes the product ID, quantity, and price at the time of the order.
type OrderItem struct {
	ProductID int
	Quantity  int
	Price     float64 // Price at the time the order was created, frozen to avoid price changes affecting past orders.
}

// OrderStatus represents the status of an order, which can be "CREATED", "COMPLETED", or "CANCELLED".
type OrderStatus string

const (
	Created   OrderStatus = "CREATED"
	Completed OrderStatus = "COMPLETED"
	Cancelled OrderStatus = "CANCELLED"
)

// Order represents a customer's order, which includes the order ID, status, items, and total price.
type Order struct {
	ID       int
	Username string
	Status   OrderStatus
	Items    []OrderItem
	Total    float64
}

// NewOrder is a constructor function that creates a new order with the given ID and items.
func NewOrder(id int, username string, items []OrderItem) (*Order, error) {
	// Validate that the order contains at least one item.
	if len(items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	return &Order{
		ID:       id,
		Username: username,
		Status:   Created,
		Items:    items,
		Total:    0.0, // Total will be calculated based on the items and their prices.
	}, nil
}

// CalculateTotal calculates the total price of the order based on the frozen prices in each OrderItem.
func (o *Order) CalculateTotal(products map[int]*product.Product) error {
	total := 0.0
	for _, item := range o.Items {
		if item.Price <= 0 {
			// Fallback to current product price if frozen price not set
			p, exists := products[item.ProductID]
			if !exists {
				return errors.New("product not found")
			}
			total += p.Price * float64(item.Quantity)
		} else {
			total += item.Price * float64(item.Quantity)
		}
	}
	o.Total = total
	return nil
}

// CompleteOrder changes the status of the order to "COMPLETED". It returns an error if the order is already completed or cancelled.
func (o *Order) CompleteOrder() error {
	if o.Status == Completed {
		return errors.New("order is already completed")
	}
	if o.Status == Cancelled {
		return errors.New("order is cancelled and cannot be completed")
	}
	o.Status = Completed
	return nil
}

// CancelOrder changes the status of the order to "CANCELLED". It returns an error if the order is already completed or cancelled.
func (o *Order) CancelOrder() error {
	if o.Status == Completed {
		return errors.New("order is already completed and cannot be cancelled")
	}
	if o.Status == Cancelled {
		return errors.New("order is already cancelled")
	}
	o.Status = Cancelled
	return nil
}
