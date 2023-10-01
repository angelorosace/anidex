package main

import (
	"anidex_api/googleApp"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// initialize Firebase APP
	ctx := context.Background()

	_, err := googleApp.BuildApp(ctx)

	if err != nil {
		fmt.Println(err)
	}
	// initialize Firebase DB
	//client, err := app.Database(ctx)
	//if err != nil {
	//	log.Fatalln("error in creating firebase DB client: ", err)
	//}

	// create ref at path user_scores/:userId
	//ref := client.NewRef("user_scores/" + fmt.Sprint(1))

	/*if err := ref.Set(context.TODO(), map[string]interface{}{"score": 40}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("score added/updated successfully!")*/

	// initialize API
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
