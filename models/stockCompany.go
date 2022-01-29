package models

import (
	"gorm.io/gorm"
)

type StockCompany struct {
	Ticker         string `gorm:"primaryKey"`
	NumberOfStocks int
}

type StockCompanySql struct {
	Sc StockCompany
	Db *gorm.DB
}

func (sc StockCompanySql) GetByTicker(ticker string) StockCompany {
	var result StockCompany

	sc.Db.Table("stock_companies").Select("ticker", "number_of_stocks").Where("ticker = ?", ticker).Scan(&result)
	return result
}

func (sc StockCompanySql) Create(payload StockCompany) {
	sc.Db.Create(&payload)
}

func (sc StockCompanySql) Update(ticker string, payload StockCompany) {
	companyModel := StockCompany{}
	sc.Db.Model(&companyModel).Where("ticker = ?", ticker).Update("number_of_stocks", payload.NumberOfStocks)
}

func (sc StockCompanySql) Delete(company StockCompany) {
	companyModel := StockCompany{}
	sc.Db.Model(&companyModel).Where("ticker = ?", company.Ticker).Delete(&company)
}