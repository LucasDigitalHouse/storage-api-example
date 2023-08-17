package dependencies

import (
	"app/cmd/server/handlers"
	"app/internal/products/storage"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-sql-driver/mysql"
)

func NewApplication(cfg *Config) *Application {
	return &Application{cfg: cfg}
}

// Application is a server application that expose Product API
var (
	ErrApplicationInternal = errors.New("internal application error")
)

type ConfigServer struct {
	Host string
	Port int
}
func (c *ConfigServer) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type Config struct {
	// database
	DbMySQL *mysql.Config
	// server
	Server  *ConfigServer
}

type Application struct {
	// config
	cfg *Config
}

func (a *Application) Run() (err error) {
	// dependencies
	// -> database
	var db *sql.DB
	db, err = sql.Open("mysql", a.cfg.DbMySQL.FormatDSN())
	if err != nil {
		err = fmt.Errorf("%w. %s", ErrApplicationInternal, err.Error())
		return
	}

	// -> products
	stProducts := storage.NewImplStorageProductMySQL(db)
	ctProducts := handlers.NewControllerProduct(stProducts)

	// -> server
	r := chi.NewRouter()

	// routes
	// -> products
	r.Get("/products/{id}", ctProducts.GetOne())
	r.Post("/products", ctProducts.Store())
	r.Put("/products/{id}", ctProducts.Update())
	r.Delete("/products/{id}", ctProducts.Delete())

	// run
	err = http.ListenAndServe(a.cfg.Server.Addr(), r)
	if err != nil {
		err = fmt.Errorf("%w. %s", ErrApplicationInternal, err.Error())
		return
	}
	
	return
}