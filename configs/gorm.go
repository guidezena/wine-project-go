package configs

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetReaderGorm() *gorm.DB {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	dialector := mysql.Open(connectionString)

	reader, err := gorm.Open(dialector)

	if err != nil {
		panic("Failed to connect database")
	}

	return reader
}

func GetWriterGorm() *gorm.DB {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	dialector := mysql.Open(connectionString)

	write, err := gorm.Open(dialector)

	if err != nil {
		panic("Failed to connect database")
	}

	return write
}

func CloseDbConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = dbSQL.Close()
	if err != nil {
		log.Fatal(err)
	}
}
