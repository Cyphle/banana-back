package main

import (
	"banana-back/db"
	"banana-back/hello"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	fmt.Println("TEST")
	fmt.Println("Accounts " + db.ToString())

	db.Add(db.BankAccount{
		Name: "My account",
	})

	fmt.Println("Accounts now: " + db.ToString())
	fmt.Println("END TEST")

	// TODO Faire des tests avec https://www.codingexplorations.com/blog/testing-gorm-with-sqlmock
	// GORM
	dsn := "host=localhost user=postgres password=postgres dbname=banana port=5432 sslmode=disable TimeZone=Europe/Paris"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	sqlDB, err := db.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.First(&product, 1)                 // find product with integer primary key
	db.First(&product, "code = ?", "D42") // find product with code D42

	fmt.Println("Record from database")
	fmt.Println(product)
	// END GORM

	// TODO HTTP Server vanilla Go et check les middleware : https://pkg.go.dev/github.com/nahojer/httprouter#section-readme

	// TODO Faire des tests echo https://echo.labstack.com/docs/testing
	// ECHO
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, hello.Hello()+", World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
	// END ECHO
}
