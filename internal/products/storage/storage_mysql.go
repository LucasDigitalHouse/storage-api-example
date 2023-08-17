package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// NewImplStorageProductMySQL returns new ImplStorageProductMySQL
func NewImplStorageProductMySQL(db *sql.DB) *ImplStorageProductMySQL {
	return &ImplStorageProductMySQL{db: db}
}

// ProductMySQL is a product model for MySQL
type ProductMySQL struct {
	ID		sql.NullInt32
	Name	sql.NullString
	Type	sql.NullString
	Count	sql.NullInt32
	Price	sql.NullFloat64
}

// ImplStorageProductMySQL is an implementation of StorageProduct interface
type ImplStorageProductMySQL struct {
	db *sql.DB
}

// GetOne returns one product by id
func (impl *ImplStorageProductMySQL) GetOne(id int) (p *Product, err error) {
	// query
	query := "SELECT id, name, type, count, price FROM products WHERE id = ?"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = impl.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}
	defer stmt.Close()

	// execute query
	row := stmt.QueryRow(id)
	if row.Err() != nil {
		err = row.Err()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			err = fmt.Errorf("%w. %v", ErrStorageProductNotFound, row.Err())
		default:
			err = fmt.Errorf("%w. %v", ErrStorageProductInternal, row.Err())
		}

		return
	}

	// scan row
	var product ProductMySQL
	err = row.Scan(&product.ID, &product.Name, &product.Type, &product.Count, &product.Price)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	// serialization
	p = new(Product)
	if product.Name.Valid {
		(*p).Name = product.Name.String
	}
	if product.Type.Valid {
		(*p).Type = product.Type.String
	}
	if product.Count.Valid {
		(*p).Count = int(product.Count.Int32)
	}
	if product.Price.Valid {
		(*p).Price = product.Price.Float64
	}

	return
}

// Store stores product
func (impl *ImplStorageProductMySQL) Store(p *Product) (err error) {
	// deserialize
	var product ProductMySQL
	if (*p).Name != "" {
		product.Name.Valid = true
		product.Name.String = (*p).Name
	}
	if (*p).Type != "" {
		product.Type.Valid = true
		product.Type.String = (*p).Type
	}
	if (*p).Count != 0 {
		product.Count.Valid = true
		product.Count.Int32 = int32((*p).Count)
	}
	if (*p).Price != 0 {
		product.Price.Valid = true
		product.Price.Float64 = (*p).Price
	}

	// query
	query := "INSERT INTO products (name, type, count, price) VALUES (?, ?, ?, ?)"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = impl.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}
	defer stmt.Close()

	// execute query
	result, err := stmt.Exec(product.Name, product.Type, product.Count, product.Price)
	if err != nil {
		errMySQL, ok := err.(*mysql.MySQLError); if ok {
			switch errMySQL.Number {
			case 1062:
				err = fmt.Errorf("%w. %v", ErrStorageProductNotUnique, err)
			default:
				err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
			}

			return
		}

		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("%w. %s", ErrStorageProductInternal, "rows affected != 1")
		return
	}

	// get last insert id
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	(*p).ID = int(lastInsertID)
	
	return
}

// Update updates product
func (impl *ImplStorageProductMySQL) Update(p *Product) (err error) {
	// deserialize
	var product ProductMySQL
	if (*p).Name != "" {
		product.Name.Valid = true
		product.Name.String = (*p).Name
	}
	if (*p).Type != "" {
		product.Type.Valid = true
		product.Type.String = (*p).Type
	}
	if (*p).Count != 0 {
		product.Count.Valid = true
		product.Count.Int32 = int32((*p).Count)
	}
	if (*p).Price != 0 {
		product.Price.Valid = true
		product.Price.Float64 = (*p).Price
	}

	// query
	query := "UPDATE products SET name = ?, type = ?, count = ?, price = ? WHERE id = ?"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = impl.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}
	defer stmt.Close()

	// execute query
	result, err := stmt.Exec(product.Name, product.Type, product.Count, product.Price, (*p).ID)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("%w. %s", ErrStorageProductInternal, "rows affected != 1")
		return
	}

	return
}

// Delete deletes product by id
func (impl *ImplStorageProductMySQL) Delete(id int) (err error) {
	// query
	query := "DELETE FROM products WHERE id = ?"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = impl.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}
	defer stmt.Close()

	// execute query
	result, err := stmt.Exec(id)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("%w. %s", ErrStorageProductInternal, "rows affected != 1")
		return
	}

	return
}