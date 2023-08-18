package storage

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestStorageMySQL_GetAll(t *testing.T) {
	// db info
	// -> table
	tableWarehouse := `
		CREATE TABLE IF NOT EXISTS warehouses (
			id INT NOT NULL AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			address VARCHAR(255) NOT NULL,
			PRIMARY KEY (id)
		);
	`
	tableProduct := `
		CREATE TABLE IF NOT EXISTS products (
			id INT NOT NULL AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			type VARCHAR(255) NOT NULL,
			count INT NOT NULL,
			price FLOAT NOT NULL,
			warehouse_id INT NOT NULL,
			PRIMARY KEY (id),
			CONSTRAINT fk_warehouse_id FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
		);
	`

	type input struct {}
	type output struct { ps []*Product; err error; errMsg string }
	type testCase struct {
		name string
		input input
		output output
		// set-up
		setUpDB func(db *sql.DB) (err error)
	}

	cases := []testCase{
		// case #1: success - no products
		{
			name: "success - no products",
			input: input{},
			output: output{
				ps: []*Product{},
				err: nil,
				errMsg: "",
			},
			setUpDB: func(db *sql.DB) (err error) {
				// create table
				_, err = db.Exec(tableWarehouse)
				if err != nil {
					err = fmt.Errorf("failed to create table. %v", err)
					return
				}
				_, err = db.Exec(tableProduct)
				if err != nil {
					err = fmt.Errorf("failed to create table. %v", err)
					return
				}
				return
			},
		},
		// case #2: success - 1 product
		{
			name: "success - 1 product",
			input: input{},
			output: output{
				ps: []*Product{
					{
						Name: "product 1",
						Type: "type 1",
						Count: 1,
						Price: 1.1,
						WarehouseId: 1,
					},
				},
				err: nil,
				errMsg: "",
			},
			setUpDB: func(db *sql.DB) (err error) {
				// create table
				_, err = db.Exec(tableWarehouse)
				if err != nil {
					err = fmt.Errorf("failed to create table. %v", err)
					return
				}
				_, err = db.Exec(tableProduct)
				if err != nil {
					err = fmt.Errorf("failed to create table. %v", err)
					return
				}

				// insert data
				query := "INSERT INTO products (name, type, count, price, warehouse_id) VALUES (?, ?, ?, ?, ?)"
				stmt, err := db.Prepare(query)
				if err != nil {
					err = fmt.Errorf("failed to prepare statement. %v", err)
					return
				}
				_, err = stmt.Exec("product 1", "type 1", 1, 1.1, 1)
				if err != nil {
					err = fmt.Errorf("failed to execute statement. %v", err)
					return
				}
				return
			},
		},
		// case #3: success - 2 products
		{
			name: "success - 2 products",
			input: input{},
			output: output{
				ps: []*Product{
					{
						Name: "product 1",
						Type: "type 1",
						Count: 1,
						Price: 1.1,
						WarehouseId: 1,
					},
					{
						Name: "product 2",
						Type: "type 2",
						Count: 2,
						Price: 2.2,
						WarehouseId: 2,
					},
				},
				err: nil,
				errMsg: "",
			},
			setUpDB: func(db *sql.DB) (err error) {
				// create table
				_, err = db.Exec(tableWarehouse)
				if err != nil {
					err = fmt.Errorf("failed to create table. %v", err)
					return
				}
				_, err = db.Exec(tableProduct)
				if err != nil {
					err = fmt.Errorf("failed to create table. %v", err)
					return
				}

				// insert data
				query := "INSERT INTO products (name, type, count, price, warehouse_id) VALUES (?, ?, ?, ?, ?)"
				stmt, err := db.Prepare(query)
				if err != nil {
					err = fmt.Errorf("failed to prepare statement. %v", err)
					return
				}
				_, err = stmt.Exec("product 1", "type 1", 1, 1.1, 1)
				if err != nil {
					err = fmt.Errorf("failed to execute statement. %v", err)
					return
				}
				_, err = stmt.Exec("product 2", "type 2", 2, 2.2, 2)
				if err != nil {
					err = fmt.Errorf("failed to execute statement. %v", err)
					return
				}
				return
			},
		},
		// case #4: error - failed to query
		{
			name: "error - failed to query",
			input: input{},
			output: output{
				ps: []*Product{},
				err: ErrStorageProductInternal,
				errMsg: "failed to query. sql: table doesn't exist",
			},
			setUpDB: func(db *sql.DB) (err error) {
				return
			},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// arrange
			// -> init db
			cfg := mysql.Config{
				User: os.Getenv("DB_MYSQL_USER"),
				Passwd: os.Getenv("DB_MYSQL_PASSWORD"),
				Net: "tcp",
				Addr: os.Getenv("DB_MYSQL_ADDR"),
				DBName: os.Getenv("DB_MYSQL_NAME"),
			}
			db, err := sql.Open("mysql", cfg.FormatDSN())
			require.NoError(t, err)
			defer db.Close()
			
			// -> init storage
			st := NewImplStorageProductMySQL(db)

			// -> init transaction
			tx, err := db.Begin()
			require.NoError(t, err)
			defer tx.Rollback()
			// -> set-up db
			err = c.setUpDB(db)
			require.NoError(t, err)

			// act
			ps, err := st.GetAll()

			// assert
			require.Equal(t, c.output.ps, ps)
			require.ErrorIs(t, err, c.output.err)
			if c.output.err != nil {
				require.EqualError(t, err, c.output.errMsg)
			}
		})
	}
}