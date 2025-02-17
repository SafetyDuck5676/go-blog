package database

import (
	"fmt"
	"log"

	"go-blog/internal/database/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open("postgres", "host=localhost user=bloguser dbname=blog sslmode=disable password=password")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	DB.AutoMigrate(&models.Post{})
	fmt.Println("Database connection established.")
}
