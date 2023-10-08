package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	animal "anidex_api/api/handlers"
	category "anidex_api/api/handlers"
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
		port = "4000"
	}

	//Animal
	if db == nil { //test without DB
		http.HandleFunc("/animal", animal.CreateAnimal)
	} else {
		http.HandleFunc("/animal", middleware.WithDatabase(db, animal.CreateAnimal))
	}

	//Category
	if db != nil {
		http.HandleFunc("categories", middleware.WithDatabase(db, category.GetCategories))
	}

	http.HandleFunc("/", getStatus)
	http.HandleFunc("/getFiles", getFiles)
	http.ListenAndServe("0.0.0.0:"+port, nil)
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
		setupRoutes(port, db)
	} else {
		setupRoutes(port, nil)
	}

}
