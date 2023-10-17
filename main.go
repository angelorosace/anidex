package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	handlers "anidex_api/api/handlers"
	"anidex_api/api/helpers"
	middleware "anidex_api/api/middleware"
	DB "anidex_api/db"
)

func getStatus(w http.ResponseWriter, r *http.Request) {
	// verify token
	authHeader := r.Header.Get("Authorization")

	// Check if the "Authorization" header is set
	if authHeader == "" {
		// Handle the case where the header is not provided
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return
	}

	_, e := helpers.VerifyToken(authHeader)
	if e != nil {
		if e.Error() == "Token is expired" {
			http.Error(w, e.Error(), http.StatusUnauthorized)
			return
		}
		// Handle the case where the header is not provided
		http.Error(w, "Invalid Signature", http.StatusUnauthorized)
		return
	}
	fmt.Println("OK")
}

func getFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	// verify token
	authHeader := r.Header.Get("Authorization")

	// Check if the "Authorization" header is set
	if authHeader == "" {
		// Handle the case where the header is not provided
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return
	}

	_, e := helpers.VerifyToken(authHeader)
	if e != nil {
		if e.Error() == "Token is expired" {
			http.Error(w, e.Error(), http.StatusUnauthorized)
			return
		}
		// Handle the case where the header is not provided
		http.Error(w, "Invalid Signature", http.StatusUnauthorized)
		return
	}

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
		http.HandleFunc("/animals/category/{category}/page/{page}", handlers.GetAnimals)

		//Login
		http.HandleFunc("/login", handlers.Login)

	} else {
		//Animal
		http.HandleFunc("/animal", middleware.WithDatabase(db, handlers.CreateAnimal))
		http.HandleFunc("/animals", middleware.WithDatabase(db, handlers.GetAnimals))

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
