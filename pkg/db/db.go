package db

import (
	"log"
	"trade/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {

	DB := Connect()
	db, err := DB.DB()
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	err = db.Ping()
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}
	log.Println("Successfully connected to the database!")

	err = DB.AutoMigrate(&models.Payload{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

func Connect() *gorm.DB {
	// dsn := "host=satao.db.elephantsql.com user=atvcirnc password=fFa7mq_RGHrfMQ1tvfsNanUhjF96sbCk dbname=atvcirnc port=5432 sslmode=disable"
	dsn := "host=localhost port=5432 user=postgres password=mysecretpassword dbname=postgres sslmode=disable"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
	}

	return DB
}
