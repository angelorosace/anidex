package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	handlers "anidex_api/api/handlers"
	middleware "anidex_api/api/middleware"
	DB "anidex_api/db"
)

func getStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("OK")
}

func getFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	entries, err := os.ReadDir(os.Getenv("RAILWAY_VOLUME_MOUNT_PATH") + "/uploaded_images")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		fmt.Println(e.Name())
	}
}

func setupRoutes(port string, db *sql.DB) {
	if port == "" {
		port = "3000"
	}

	if db == nil { //test without DB

		//Animal
		http.HandleFunc("/animal", handlers.CreateAnimal)
		http.HandleFunc("/animals/category/{category}/page/{page}", handlers.GetAnimalsHandler)

		//Login
		http.HandleFunc("/login", handlers.Login)

	} else {
		//Animal
		http.HandleFunc("/animal", middleware.WithDatabase(db, handlers.CreateAnimal))
		http.HandleFunc("/animals", middleware.WithDatabase(db, handlers.GetAnimalsHandler))

		//Images
		http.HandleFunc("/images", handlers.GetImageByPath)

		//Category
		http.HandleFunc("/categories", middleware.WithDatabase(db, handlers.GetCategories))

		//Stats
		http.HandleFunc("/stats", middleware.WithDatabase(db, handlers.GetStats))

		//Login
		http.HandleFunc("/login", middleware.WithDatabase(db, handlers.Login))
	}

	http.HandleFunc("/", getStatus)
	http.HandleFunc("/getFiles", getFiles)
	err := http.ListenAndServe("0.0.0.0:"+port, nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func main() {

	port := os.Getenv("PORT")

	// Initialize DB
	if port != "" {
		db, err := DB.InitializeDB()
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Connection with DB established!")
		}
		defer db.Close()
		fmt.Println("Server online reachable at port", port)
		setupRoutes(port, db)
	} else {
		fmt.Println("Server online reachable at port 3000")
		setupRoutes(port, nil)
	}

}
