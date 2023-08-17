package storage

import "errors"

// Product is a product model
type Product struct {
	ID		int
	Name    string
	Type	string
	Count	int
	Price	float64
}

// StorageProduct is an interface for product storage
type StorageProduct interface {
	// GetOne returns one product by id
	GetOne(id int) (p *Product, err error)

	// Store stores product
	Store(p *Product) (err error)

	// Update updates product
	Update(p *Product) (err error)

	// Delete deletes product by id
	Delete(id int) (err error)
}

var (
	ErrStorageProductInternal = errors.New("internal storage product error")
	ErrStorageProductNotFound = errors.New("storage product not found")
	ErrStorageProductNotUnique = errors.New("storage product not unique")
)