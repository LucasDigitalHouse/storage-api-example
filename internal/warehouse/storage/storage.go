package storage

type WarehouseAttributes struct {
	Name    string
	Address string
}

// Warehouse is a warehouse model
type Warehouse struct {
	// ID is a unique identifier
	ID int

	// Attr is the warehouse attributes
	Attr *WarehouseAttributes
}
