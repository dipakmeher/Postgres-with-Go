package models

import "gorm.io/gorm"

/*
* gorm keyword allows us to define primary key, etc.
* uint will autoincrement when interacting with API
* Rest others are pointers, referring to the same variables
 */
type Books struct {
	ID        uint    `gorm:"primary key; autoIncrement" json:"id"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:"publiser"`
}

/*
* MongoDB will create database if it not there
* Postgres needs to have that database existing
* MigrateBook will help creating database in Postgres if it does not exist
 */
func MigrateBook(db *gorm.DB) error {
	err := db.AutoMigrate(&Books{})
	return err
}
