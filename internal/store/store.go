package store

import (
	"fmt"

	"github.com/DiegoAmin/AmazonClone_PAP/internal/logger"
	"github.com/DiegoAmin/AmazonClone_PAP/internal/order"
	"github.com/DiegoAmin/AmazonClone_PAP/internal/product"
)

// Store represents the main structure of the store, which contains a map of products and a map of orders.
type Store struct {
	Products map[int]*product.Product
	Orders   map[int]*order.Order
}

// NewStore is a constructor function that initializes and returns a new Store instance with empty product and order maps.
func NewStore() (*Store, error) {
	if err := logger.Init("store.log"); err != nil {
		return nil, err
	}
	return &Store{
		Products: make(map[int]*product.Product),
		Orders:   make(map[int]*order.Order),
	}, nil
}

// AddProduct adds a new product to the store. It returns an error if a product with the same ID already exists.
func (s *Store) AddProduct(p product.Product) error {
	if _, exists := s.Products[p.ID]; exists {
		logger.Log(fmt.Sprintf("ERROR: product with ID %d already exists", p.ID))
		return fmt.Errorf("product with ID %d already exists", p.ID)
	}
	s.Products[p.ID] = &p
	logger.Log(fmt.Sprintf("ADMIN: product added: ID=%d, Name=%s, Price=%.2f, Stock=%d", p.ID, p.Name, p.Price, p.Stock))
	return nil
}

// CreateOrder creates a new order with the given items. It validates stock availability before making any changes.
// Prices are frozen at the time of order creation to avoid price changes affecting past orders.
func (s *Store) CreateOrder(items []order.OrderItem, username string) (*order.Order, error) {
	orderID := len(s.Orders) + 1

	// Step 1: Validate that all products exist and have enough stock before making any changes.
	for _, item := range items {
		product, exists := s.Products[item.ProductID]
		if !exists {
			logger.Log(fmt.Sprintf("ERROR: product with ID %d not found", item.ProductID))
			return nil, fmt.Errorf("product with ID %d not found", item.ProductID)
		}
		if product.Stock < item.Quantity {
			logger.Log(fmt.Sprintf("ERROR: not enough stock for product ID %d", item.ProductID))
			return nil, fmt.Errorf("not enough stock for product ID %d", item.ProductID)
		}
	}

	// Step 2: Freeze prices at the time of order creation.
	frozenItems := make([]order.OrderItem, len(items))
	for i, item := range items {
		frozenItems[i] = order.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     s.Products[item.ProductID].Price, // freeze current price
		}
	}

	newOrder, err := order.NewOrder(orderID, username, frozenItems)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: failed to create order: %s", err.Error()))
		return nil, err
	}

	// Step 3: Only if everything is fine, reduce the stock.
	for _, item := range items {
		s.Products[item.ProductID].Stock -= item.Quantity
	}

	// Step 4: Calculate the total price of the order using frozen prices.
	err = newOrder.CalculateTotal(s.Products)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: failed to calculate total: %s", err.Error()))
		return nil, err
	}

	s.Orders[orderID] = newOrder
	logger.Log(fmt.Sprintf("CLIENT: order created: ID=%d, Total=%.2f", newOrder.ID, newOrder.Total))
	return newOrder, nil
}

// CompleteOrder marks an order as completed. It returns an error if the order does not exist or if it is already completed or cancelled.
func (s *Store) CompleteOrder(orderID int) error {
	order, exists := s.Orders[orderID]
	if !exists {
		logger.Log(fmt.Sprintf("ERROR: order with ID %d not found", orderID))
		return fmt.Errorf("order with ID %d not found", orderID)
	}
	if err := order.CompleteOrder(); err != nil {
		logger.Log(fmt.Sprintf("ERROR: failed to complete order ID %d: %s", orderID, err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("ADMIN: order completed: ID=%d", orderID))
	return nil
}

// CancelOrder marks an order as cancelled and restores the stock of all items in the order.
// It returns an error if the order does not exist or if it is already completed or cancelled.
func (s *Store) CancelOrder(orderID int, username string) error {
	o, exists := s.Orders[orderID]
	if !exists {
		logger.Log(fmt.Sprintf("ERROR: order with ID %d not found", orderID))
		return fmt.Errorf("order with ID %d not found", orderID)
	}
	if o.Username != username {
		logger.Log(fmt.Sprintf("ERROR: user %s is not authorized to cancel order ID %d", username, orderID))
		return fmt.Errorf("user %s is not authorized to cancel this order", username)
	}
	if err := o.CancelOrder(); err != nil {
		logger.Log(fmt.Sprintf("ERROR: failed to cancel order ID %d: %s", orderID, err.Error()))
		return err
	}

	// Restore the stock for each item in the cancelled order.
	for _, item := range o.Items {
		if product, exists := s.Products[item.ProductID]; exists {
			product.Stock += item.Quantity
		}
	}

	logger.Log(fmt.Sprintf("CLIENT: order cancelled: ID=%d, stock restored", orderID))
	return nil
}

// UpdatePrice updates the price of a product. It returns an error if the product does not exist or if the new price is negative.
func (s *Store) UpdatePrice(productID int, price float64) error {
	product, exists := s.Products[productID]
	if !exists {
		logger.Log(fmt.Sprintf("ERROR: product with ID %d not found", productID))
		return fmt.Errorf("product with ID %d not found", productID)
	}
	if err := product.UpdatePrice(price); err != nil {
		logger.Log(fmt.Sprintf("ERROR: failed to update price for product ID %d: %s", productID, err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("ADMIN: product price updated: ID=%d, NewPrice=%.2f", productID, price))
	return nil
}

// UpdateStock updates the stock of a product. It returns an error if the product does not exist or if the new stock is negative.
func (s *Store) UpdateStock(productID int, stock int) error {
	product, exists := s.Products[productID]
	if !exists {
		logger.Log(fmt.Sprintf("ERROR: product with ID %d not found", productID))
		return fmt.Errorf("product with ID %d not found", productID)
	}
	if err := product.UpdateStock(stock); err != nil {
		logger.Log(fmt.Sprintf("ERROR: failed to update stock for product ID %d: %s", productID, err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("ADMIN: product stock updated: ID=%d, NewStock=%d", productID, stock))
	return nil
}

// OrderHistory returns all orders in the store as a slice.
func (s *Store) OrderHistory() []*order.Order {
	logger.Log("ADMIN: order history requested")
	orders := make([]*order.Order, 0, len(s.Orders))
	for _, order := range s.Orders {
		orders = append(orders, order)
	}
	return orders
}

// OrderHistoryByUser returns all orders for a specific user as a slice. It returns an error if no orders are found for the user.
func (s *Store) OrderHistoryByUser(username string) []*order.Order {
	logger.Log(fmt.Sprintf("CLIENT: order history requested for user: %s", username))
	orders := make([]*order.Order, 0)
	for _, order := range s.Orders {
		if order.Username == username {
			orders = append(orders, order)
		}
	}
	return orders
}

// ListProducts returns all products in the store as a slice.
func (s *Store) ListProducts() []*product.Product {
	logger.Log("CLIENT: product list requested")
	products := make([]*product.Product, 0, len(s.Products))
	for _, product := range s.Products {
		products = append(products, product)
	}
	return products
}

// GetProduct returns a product by its ID. It returns the product and a boolean indicating if the product was found.
func (s *Store) GetProduct(id int) (*product.Product, bool) {
	p, exists := s.Products[id]
	return p, exists
}
