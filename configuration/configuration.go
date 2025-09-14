package configuration

import (
	"log"
	"os"

	"fmt"

	"github.com/U-T-kuroitigo/Saikyo_UI/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Configuration creates a struct for the environment variables
type Configuration struct {
	Server   string
	Port     string
	User     string
	Password string
	Database string
}

// GetConfiguration gets the configuration from the environment variables
func GetConfiguration() Configuration {
	var c Configuration
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// 環境変数を取得
	c.Server = os.Getenv("Server")
	c.Port = os.Getenv("Port")
	c.User = os.Getenv("User")
	c.Password = os.Getenv("Password")
	c.Database = os.Getenv("Database")

	return c
}

// InitDB initializes the database connection and performs auto-migration
func InitDB() *gorm.DB {
	c := GetConfiguration()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", c.User, c.Password, c.Server, c.Port, c.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// AutoMigrate the models
	db.AutoMigrate(&models.User{})

	return db
}
