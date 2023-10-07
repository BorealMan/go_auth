package database

import (
	"app/config"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	DB *sqlx.DB
)

func InitDB() {

	var err error
	DB, err = sqlx.Connect("mysql", GenerateDBURL())

	if err != nil {
		log.Fatalf(`Failed To Connect To Database: %v+\n`, err)
	}

	fmt.Println("Successfully Connected To Database")
	Seed()
}

func GenerateDBURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=EST",
		config.DB_USERNAME,
		config.DB_PASSWORD,
		config.DB_HOST,
		config.DB_PORT,
		config.DB_DATABASE,
	)
}
