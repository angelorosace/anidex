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

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()

	if db == nil { //test without DB

		//Animal
		r.HandleFunc("/animal", handlers.CreateAnimal).Methods("GET")
		r.HandleFunc("/animals/category/{category}/page/{page}", handlers.GetAnimalsByCategory).Methods("GET")

	} else {
		//Animal
		r.HandleFunc("/animal", middleware.WithDatabase(db, handlers.CreateAnimal)).Methods("POST")
		r.HandleFunc("/animals/category/{category}/page/{page}", middleware.WithDatabase(db, handlers.GetAnimalsByCategory)).Methods("GET")
		r.HandleFunc("/animals/id/{id}", middleware.WithDatabase(db, handlers.GetAnimalById)).Methods("GET")

		//Images
		r.HandleFunc("/images/photo/{photo}", handlers.GetImageByPath).Methods("GET")

		//Category
		r.HandleFunc("/categories", middleware.WithDatabase(db, handlers.GetCategories)).Methods("GET")

		//Stats
		r.HandleFunc("/stats", middleware.WithDatabase(db, handlers.GetStats)).Methods("GET")
	}

	http.Handle("/", r)
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
