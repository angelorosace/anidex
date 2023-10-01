package main

import (
	"fmt"
	"log"
	"os"

	DB "anidex_api/db"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// Initialize DB
	db, err := DB.InitializeDB()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection with DB established!")
	}

	db.Exec("CREATE TABLE Test (ID int);")

	// initialize API Service
	fiberApp := fiber.New()

	fiberApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	log.Fatal(fiberApp.Listen("0.0.0.0:" + port))

}
