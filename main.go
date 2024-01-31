package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("Failed to connect to database")
	}

	db.AutoMigrate(&Book{})

	app := fiber.New()

	app.Get("/books", func(c *fiber.Ctx) error {
		return getBooks(db, c)
	})
	app.Get("/books/:id", func(c *fiber.Ctx) error {
		return getBook(db, c)
	})
	app.Post("/books", func(c *fiber.Ctx) error {
		return createBook(db, c)
	})
	app.Put("/books/:id", func(c *fiber.Ctx) error {
		return updateBook(db, c)
	})
	app.Delete("/books/:id", func(c *fiber.Ctx) error {
		return deleteBook(db, c)
	})

	log.Fatal(app.Listen(":8000"))
}
