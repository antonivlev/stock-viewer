/*
Package database provides functions for saving and retrieving objects from the database.
TODO: add tests
*/
package database

import (
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

var gormDB *gorm.DB

type Search struct {
	SearchTime time.Time
	Stock      string
}

type stockData struct {
	Stock string `gorm:"primary_key"`
	Data  string
}

// Connects the database
func SetupDatabase() error {
	db, errConnect := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=stockapp password=12345 sslmode=disable")
	if errConnect != nil {
		return errConnect
	}

	db.Exec("delete from searches;delete from stock_data;")

	// create table for searches
	db.AutoMigrate(&Search{})
	db.AutoMigrate(&stockData{})

	gormDB = db
	return nil
}

func IsStockCached(stock string) bool {
	var stockData stockData
	stock = strings.ToUpper(stock)
	gormDB.Where("stock = ?", stock).Find(&stockData)
	return stockData.Stock != ""
}

func GetCachedStockData(stock string) []byte {
	var stockData stockData
	stock = strings.ToUpper(stock)
	gormDB.Where("stock = ?", stock).Find(&stockData)

	return []byte(stockData.Data)
}

func SaveStockData(stock string, stockDataBytes []byte) {
	stock = strings.ToUpper(stock)
	stockData := stockData{
		Stock: stock,
		Data:  string(stockDataBytes),
	}
	gormDB.Save(&stockData)
}

func SaveSearch(s Search) error {
	if gormDB == nil {
		return errors.New("database not configured, check SetupDatabase()")
	}

	s.Stock = strings.ToUpper(s.Stock)

	// TODO: check for save errors
	gormDB.Save(&s)
	return nil
}

func GetSearches() ([]Search, error) {
	if gormDB == nil {
		return nil, errors.New("database not configured, check SetupDatabase()")
	}

	var searches []Search
	gormDB.Find(&searches)
	return searches, nil
}
