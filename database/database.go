/*
Package database provides functions for saving and retrieving objects from the database.
TODO: add tests
*/
package database

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

var gormDB *gorm.DB

type Search struct {
	SearchTime time.Time
	Stock     string
}

// Connects the database
func SetupDatabase() error {
	db, errConnect := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=stockapp password=12345 sslmode=disable")
	if errConnect != nil {
		return errConnect
	}

	db.Exec("delete from searches;")

	// create table for searches
	db.AutoMigrate(&Search{})

	gormDB = db
	return nil
}

func SaveSearch(s Search) error {
	if gormDB == nil {
		return errors.New("database not configured, check SetupDatabase()")
	}

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
