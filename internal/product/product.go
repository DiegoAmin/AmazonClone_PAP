package product

import (
	"errors"
	"regexp"
)

// Product represents a product in the inventory.

type Product struct {
	//All of the atributes needs to be public to be able to be used outside of the package.
	ID    int
	Name  string
	Price float64
	Stock int
}

// Constructor function to create a new product. It validates the price and stock before creating the product.
func NewProduct(id int, name string, price float64, stock int) (*Product, error) {
	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	// Validate that the name is alphanumeric using a regular expression.

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, name)

	if !matched {
		return nil, errors.New("name must be alphanumeric")
	}

	return &Product{
		ID:    id,
		Name:  name,
		Price: price,
		Stock: stock,
	}, nil
}

// UpdatePrice updates the price of the product. It validates the new price before updating it.
func (p *Product) UpdatePrice(price float64) error {
	if price < 0 {
		return errors.New("price cannot be negative")
	}
	p.Price = price
	return nil
}

// UpdateStock updates the stock of the product. It validates the new stock before updating it.
func (p *Product) UpdateStock(stock int) error {
	if stock < 0 {
		return errors.New("stock cannot be negative")
	}
	p.Stock = stock
	return nil
}
