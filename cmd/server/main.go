package main

import (
	"app/cmd/server/dependencies"
	"os"

	"github.com/go-sql-driver/mysql"
)

func main() {
	// env
	// ...


	// app
	// -> cfg
	cfg := &dependencies.Config{
		// database
		DbMySQL: &mysql.Config{
			User: os.Getenv("DB_MYSQL_USER"),
			Passwd: os.Getenv("DB_MYSQL_PASSWORD"),
			Net: "tcp",
			Addr: os.Getenv("DB_MYSQL_ADDR"),
			DBName: os.Getenv("DB_MYSQL_DATABASE"),
			ParseTime: true,
		},
		// server
		Server: &dependencies.ConfigServer{
			Host: os.Getenv("SERVER_HOST"),
			Port: 8080,
		},
	}

	app := dependencies.NewApplication(cfg)

	// -> run
	if err := app.Run(); err != nil {
		panic(err)
	}
}