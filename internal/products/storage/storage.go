package storage

import (
	storageWh "app/internal/warehouse/storage"
	"errors"
)

// Product is a product model
type Product struct {
	ID			int
	Name    	string
	Type		string
	Count		int
	Price		float64
	WarehouseId int
}

type ProductWarehouse struct {
	// Product
	ID			int
	Name    	string
	Type		string
	Count		int
	Price		float64
	// Warehouse Attributes
	WarehouseAttr *storageWh.WarehouseAttributes
}

// StorageProduct is an interface for product storage
type StorageProduct interface {
	// GetOne returns one product by id
	GetOne(id int) (p *Product, err error)

	// GetAll returns all products
	GetAll() (ps []*Product, err error)

	// GetOneWithWarehouse returns one product by id with warehouse info
	GetOneWithWarehouse(id int) (p *ProductWarehouse, err error)

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
	ErrStorageProductRelation = errors.New("storage product relation error")
)