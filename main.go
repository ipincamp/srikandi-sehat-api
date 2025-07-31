package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// Database connection function
func connectDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := connectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()
	log.Println("✅ Database connected successfully!")

	app := fiber.New() // Development
	/*
		// Production
		app := fiber.New(fiber.Config{
			Prefork: true,
		})
	*/

	app.Get("/", func(c *fiber.Ctx) error {
		if err := db.Ping(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "database connection is lost",
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Hello, World! API is running and database is connected.",
		})
	})

	port := os.Getenv("API_PORT")
	log.Printf("🚀 Server is running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
