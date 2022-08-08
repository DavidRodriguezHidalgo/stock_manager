package main

import (
	"fmt"
	"strconv"

	"./models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Tview
var pages = tview.NewPages()
var companyText = tview.NewTextView()
var app = tview.NewApplication()
var form = tview.NewForm()
var stockCompaniesList = tview.NewList().ShowSecondaryText(false)
var flex = tview.NewFlex()
var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("(a) to add a new stock company \n(d) to delete a stock company \n(q) to quit")

var stockCompanies = make([]models.StockCompany, 0)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.StockCompany{})
	loadCompanies(db)
	reloadCompanies(db)

	stockCompaniesList.SetSelectedFunc(func(index int, name string, second_name string, shortcut rune) {
		setConcatText(&stockCompanies[index])
	})

	fmt.Println(len(stockCompanies))
	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(stockCompaniesList, 0, 1, true).
			AddItem(companyText, 0, 4, false), 0, 6, false).
		AddItem(text, 0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 {
			app.Stop()
		} else if event.Rune() == 97 {
			form.Clear(true)
			addStockCompanyForm(db)
			pages.SwitchToPage("Add Stock Company")
		} else if event.Rune() == 100 {
			form.Clear(true)
			deleteStockCompanyForm(db)
			pages.SwitchToPage("Delete Stock Company")
		}

		return event
	})

	pages.AddPage("Menu", flex, true, true)
	pages.AddPage("Add Stock Company", form, true, false)
	pages.AddPage("Delete Stock Company", form, true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func addStockCompanyForm(db *gorm.DB) *tview.Form {
	companyModel := models.StockCompany{}
	company := models.StockCompanySql{Db: db, Sc: companyModel}
	company1 := models.StockCompany{}

	form.AddInputField("Ticker", "", 20, nil, func(ticker string) {
		company1.Ticker = ticker
	})

	form.AddInputField("Number of stocks", "", 20, nil, func(number string) {
		num, _ := strconv.Atoi(number)
		company1.NumberOfStocks = num
	})

	form.AddButton("Save", func() {
		company.Create(company1)
		stockCompanies = append(stockCompanies, models.StockCompany{Ticker: company1.Ticker, NumberOfStocks: company1.NumberOfStocks})
		reloadCompanies(db)
		pages.SwitchToPage("Menu")
	})

	return form
}

func deleteStockCompanyForm(db *gorm.DB) *tview.Form {
	companyModel := models.StockCompany{}
	company := models.StockCompanySql{Db: db, Sc: companyModel}
	var tickerToSearch = ""

	form.AddInputField("Ticker", "", 20, nil, func(ticker string) {
		tickerToSearch = ticker
	})

	form.AddButton("Save", func() {
		company.Delete(company.GetByTicker(tickerToSearch))
		reloadCompanies(db)
		pages.SwitchToPage("Menu")
	})

	return form
}

func loadCompanies(db *gorm.DB) {
	stockCompanies = make([]models.StockCompany, 0)
	companyModel := models.StockCompany{}
	company := models.StockCompanySql{Db: db, Sc: companyModel}
	companies := company.GetAll()

	for _, c := range companies {
		stockCompanies = append(stockCompanies, c)
	}
}

func reloadCompanies(db *gorm.DB) {
	stockCompaniesList.Clear()

	loadCompanies(db)
	for index, company := range stockCompanies {
		stockCompaniesList.AddItem(company.Ticker+" "+strconv.Itoa(company.NumberOfStocks), " ", rune(49+index), nil)
	}
}

func setConcatText(company *models.StockCompany) {
	companyText.Clear()
	text := company.Ticker + " " + strconv.Itoa(company.NumberOfStocks) + "\n"
	companyText.SetText(text)
}
