package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db           *gorm.DB
	dbClientOnce sync.Once
)

func dsn() string {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load db configuration from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_DATABASE")

	// Validasi variabel environment penting
	if host == "" || user == "" || pass == "" || db_name == "" || port == "" {
		log.Fatal("❌ Missing required database environment variables (DB_HOST, DB_USERNAME, DB_PASSWORD, DB_DATABASE, DB_PORT)")
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, db_name, port,
	)
}

func ConnectDB() *gorm.DB {
	dbClientOnce.Do(func() {
		// Inisialisasi koneksi GORM
		var err error
		db, err = gorm.Open(postgres.Open(dsn()), &gorm.Config{
			Logger:         logger.Default.LogMode(logger.Info),
			TranslateError: true,
		})

		if err != nil {
			log.Fatalf("❌ Failed to connect to database: %v", err)
		}
	})

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get *sql.DB: %v", err)
		return nil
	}

	// ping connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatalln("Failed to ping database: ", err)
		return nil
	}

	log.Println("Successfully connect database.")

	return db
}
