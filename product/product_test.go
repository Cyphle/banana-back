package product

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})

	// CASE 1
	createProduct(db)

	// CASE 2
	// fmt.Println(mock)
	rows := sqlmock.NewRows([]string{"Code", "Price"}).AddRow("D43", 100)
	mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
	createProduct(db)
}
