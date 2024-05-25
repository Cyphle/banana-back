package product

// GORM TEST EXAMPLE
//type Product struct {
//	gorm.Model
//	Code  string
//	Price uint
//}
//
//func createProduct(db *gorm.DB) {
//	var product Product
//	err := db.First(&product, "code = ?", "D42").Error
//	if err != nil {
//		return
//	}
//	db.Create(&Product{Code: "D42", Price: 100})
//}
