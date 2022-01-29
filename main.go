package main

import (
	"fmt"

	"./models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.StockCompany{})

	
	companyModel := models.StockCompany{}
	company := models.StockCompanySql{Db: db, Sc: companyModel}
	company1 := models.StockCompany{Ticker: "ZOT", NumberOfStocks: 100}
	company.Create(company1)
	fmt.Println(company.GetByTicker("ZOT"))
	company.Update("ZOT", models.StockCompany{NumberOfStocks: 200})
	fmt.Println(company.GetByTicker("ZOT"))
	//company.Delete(company.GetByTicker("ZOT"))
}
