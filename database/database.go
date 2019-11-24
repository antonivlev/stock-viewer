/*
Package database provides functions for saving and retrieving objects from the database.
TODO: add tests
*/
package database

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	gormDB           *gorm.DB
	deletionInterval time.Duration
)

type Search struct {
	SearchTime time.Time
	Stock      string
}

type stockData struct {
	Stock          string `gorm:"primary_key"`
	Data           string
	ExpirationTime time.Time
}

// Connects the database
func SetupDatabase(configFile string) error {
	var config map[string]string
	configBytes, errConfigRead := ioutil.ReadFile(configFile)
	if errConfigRead != nil {
		return errConfigRead
	}
	errParseConfig := json.Unmarshal(configBytes, &config)
	if errParseConfig != nil {
		return errParseConfig
	}

	// connect to db
	db, errConnect := gorm.Open("postgres",
		"host="+config["hostname"]+
			" port="+config["port"]+
			" user="+config["user"]+
			" dbname="+config["dbname"]+
			" sslmode=disable",
	)
	if errConnect != nil {
		return errConnect
	}

	numSeconds, errParseInterval := strconv.Atoi(config["cacheExpirationInterval"])
	if errParseInterval != nil {
		return errParseInterval
	}
	deletionInterval = time.Second * time.Duration(numSeconds)

	// create table for searches
	db.AutoMigrate(&Search{})
	db.AutoMigrate(&stockData{})
	db.Exec("delete from searches;delete from stock_data;")

	gormDB = db
	go deleteExpiredData()
	return nil
}

func deleteExpiredData() {
	for {
		time.Sleep(time.Second * 10)
		log.Println("checking expired data...")
		var expiredOnes []stockData
		gormDB.Where("expiration_time < ?", time.Now()).Find(&expiredOnes)
		for _, expiredOne := range expiredOnes {
			log.Printf("%s expired at %s, deleting\n", expiredOne.Stock, expiredOne.ExpirationTime)
			gormDB.Delete(&expiredOne)
		}
	}
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
		Stock:          stock,
		Data:           string(stockDataBytes),
		ExpirationTime: time.Now().Add(deletionInterval),
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
