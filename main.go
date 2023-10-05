package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	DB "anidex_api/db"
)

type animalPostRequest struct {
	Photo []struct {
	} `json:"photo"`
	Name        string   `json:"name"`
	Taxonomy    string   `json:"taxonomy"`
	Etymology   string   `json:"etymology"`
	Iucn        []string `json:"iucn"`
	Geo         string   `json:"geo"`
	Migration   string   `json:"migration"`
	Habitat     string   `json:"habitat"`
	Dimensions  string   `json:"dimensions"`
	Ds          string   `json:"ds"`
	Diet        string   `json:"diet"`
	Description string   `json:"description"`
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("OK")
}

func parseAnimalRequest(m *multipart.Form) {

	for _, v := range m.File["photo[]"] {
		fmt.Println(v.Filename, ":", v.Size)
		file, err := v.Open()
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		/*
			// Create directory
			//dirPath := "/Users/accilo/Desktop/angelo/anidex_api/temp-img"
			dirPath := os.Getenv("RAILWAY_VOLUME_MOUNT_PATH") + "/uploaded_images"
			if _, err := os.Stat(dirPath); errors.Is(err, os.ErrNotExist) {
				err := os.Mkdir(dirPath, os.ModePerm)
				if err != nil {
					fmt.Println(err)
				}
			}

			// Create file
			dst, err := os.Create(filepath.Join(dirPath, filepath.Base(v.Filename)))
			if err != nil {
				fmt.Println(err)
				return
			}
			defer dst.Close()
			if _, err = io.Copy(dst, file); err != nil {
				fmt.Println(err)
				return
			}*/
	}
	fmt.Println("Successfully Uploaded Images")
}

func postAnimal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	r.ParseMultipartForm(10 << 20)

	parseAnimalRequest(r.MultipartForm)

}

func setupRoutes(port string) {
	if port == "" {
		port = "4000"
	}

	http.HandleFunc("/", getStatus)
	http.HandleFunc("/animal", postAnimal)
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
	}

	setupRoutes(port)

}
