package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"anidex_api/api/helpers"
	animal "anidex_api/domain/animal"
	responses "anidex_api/http/responses"
)

type AnimalRequest struct {
	Photos      string `json:"photo"`
	Category    string `json:"category"`
	Name        string `json:"name"`
	Taxonomy    string `json:"taxonomy"`
	Etymology   string `json:"etymology"`
	Iucn        string `json:"iucn"`
	Geo         string `json:"geo"`
	Migration   string `json:"migration"`
	Habitat     string `json:"habitat"`
	Dimensions  string `json:"dimensions"`
	Ds          string `json:"ds"`
	Diet        string `json:"diet"`
	Description string `json:"description"`
}

type Animal struct {
	ID          int    `json:"id"`
	Photos      string `json:"photo"`
	Category    string `json:"category"`
	Name        string `json:"name"`
	Taxonomy    string `json:"taxonomy"`
	Etymology   string `json:"etymology"`
	Iucn        string `json:"iucn"`
	Geo         string `json:"geo"`
	Migration   string `json:"migration"`
	Habitat     string `json:"habitat"`
	Dimensions  string `json:"dimensions"`
	Ds          string `json:"ds"`
	Diet        string `json:"diet"`
	Description string `json:"description"`
}

type AnimalPreview struct {
	ID     int    `json:"id"`
	Photos string `json:"photos"`
	Name   string `json:"name"`
}

type AnimalResponse struct {
	AnimalData AnimalRequest `json:"animalData"`
	Error      string        `json:"error"`
	Message    string        `json:"message"`
	Status     int           `json:"status"`
}

func getDataFromMap(key string, originData map[string][]string) ([]string, error) {
	if data, exists := originData[key]; exists {
		return data, nil
	}
	return nil, fmt.Errorf("%s is not present in provided data", key)
}

func (ar *AnimalRequest) readAnimalRequestValues(values map[string][]string) error {
	for _, v := range animal.ANIMAL_POST_REQUEST_MANDATORY_FIELDS {
		if v == "photo[]" {
			continue
		}

		var structFieldName string
		if v == "iucn[]" {
			structFieldName = strings.ToUpper(v[0:1]) + v[1:]
			structFieldName = strings.ReplaceAll(structFieldName, "[]", "")
		} else {
			structFieldName = strings.ToUpper(v[0:1]) + v[1:]
		}

		value, err := getDataFromMap(v, values)

		if v == "description" && err != nil {
			ar.Description = ""
			continue
		}

		if v == "iucn[]" && err == nil {
			ar.Iucn = strings.Join(value, ",")
			continue
		}

		if err != nil {
			return err
		}

		reflect.ValueOf(ar).Elem().FieldByName(structFieldName).SetString(value[0])
	}

	return nil
}

func animalRequestErrorResponse(w http.ResponseWriter, err error) []byte {
	response := AnimalResponse{
		Error:   err.Error(),
		Status:  http.StatusBadRequest,
		Message: "",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return jsonResponse
}

func animalRequestSuccessResponse(w http.ResponseWriter, a AnimalRequest) []byte {
	response := AnimalResponse{
		AnimalData: a,
		Status:     http.StatusCreated,
		Message:    "New animal species successfully registered in the Anidex",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return jsonResponse
}

func storePhotosAndCollectPaths(m *multipart.Form) ([]string, int32, []error) {
	var photoPaths []string
	var storedPhotosCount int32
	var errs []error

	for _, v := range m.File["photo[]"] {
		file, err := v.Open()
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		dirPath := os.Getenv("RAILWAY_VOLUME_MOUNT_PATH") + "/uploaded_images"
		if _, err := os.Stat(dirPath); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(dirPath, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				errs = append(errs, err)
				continue
			}
		}

		newFilePath := filepath.Join(dirPath, filepath.Base(v.Filename))

		// Create file
		dst, err := os.Create(newFilePath)
		if err != nil {
			fmt.Println(err)
			errs = append(errs, err)
			continue
		}
		defer dst.Close()

		if _, err = io.Copy(dst, file); err != nil {
			fmt.Println(err)
			errs = append(errs, err)
			continue
		}

		photoPaths = append(photoPaths, newFilePath)
		storedPhotosCount += 1
	}
	fmt.Println("Successfully Uploaded Image")
	return photoPaths, storedPhotosCount, errs
}

func (ar *AnimalRequest) buildAnimalRequest(m *multipart.Form) error {

	photoPaths, uploadedPhotosCount, errs := storePhotosAndCollectPaths(m)

	if len(errs) > 0 {
		return fmt.Errorf("The upload of photos produced an error: %v", errs)
	}

	if len(photoPaths) != int(uploadedPhotosCount) {
		return fmt.Errorf("The upload of photos produced an error: %v", errs)
	}

	ar.Photos = strings.Join(photoPaths, ",")

	err := ar.readAnimalRequestValues(m.Value)
	if err != nil {
		return err
	}

	return nil

}

func CDAnimal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // Add DELETE to the allowed methods
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	// Check if it's an OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method == http.MethodDelete {
		DeleteAnimal(w, r)
	}
	if r.Method == http.MethodPost {
		CreateAnimal(w, r)
	}
}

func DeleteAnimal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")

	idToDelete := r.URL.Query().Get("id")

	if idToDelete == "" {
		resp, err := responses.MissingURLParametersResponse(w)
		if err != nil {
			http.Error(w, "", http.StatusBadGateway)
		}
		w.Write(resp)
		return
	}

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
		responses.CustomResponse(w, nil, e.Error(), http.StatusUnauthorized, e.Error())
		return
	}

	//retrieve DB from context
	db := r.Context().Value("db").(*sql.DB)

	//save data in mysql
	stmt, err := db.Prepare("DELETE FROM animals WHERE id=" + idToDelete)
	if err != nil {
		res, err := responses.CustomResponse(w, nil, "myslq produced an error", http.StatusInternalServerError, err.Error())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(res)
		return
	}
	defer stmt.Close()

	res, err := responses.CustomResponse(w, nil, "Animal Succesfully deleted", http.StatusOK, "")
	w.Write(res)
}

func CreateAnimal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
	w.Header().Set("Content-Type", "application/json")
	r.ParseMultipartForm(10 << 20)

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
		responses.CustomResponse(w, nil, e.Error(), http.StatusUnauthorized, e.Error())
		return
	}

	var animalRequest AnimalRequest
	err := animalRequest.buildAnimalRequest(r.MultipartForm)
	if err != nil {
		w.Write(animalRequestErrorResponse(w, err))
		return
	}

	//retrieve DB from context
	db := r.Context().Value("db").(*sql.DB)

	//save data in mysql
	stmt, err := db.Prepare("INSERT INTO animals (photos,name,taxonomy,etymology,iucn,geo,migration,habitat,dimensions,ds,diet,description,category) VALUES (?, ?,?, ?,?, ?,?, ?,?, ?,?, ?,?)")
	if err != nil {
		w.Write(animalRequestErrorResponse(w, err))
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		animalRequest.Photos,
		animalRequest.Name,
		animalRequest.Taxonomy,
		animalRequest.Etymology,
		animalRequest.Iucn,
		animalRequest.Geo,
		animalRequest.Migration,
		animalRequest.Habitat,
		animalRequest.Dimensions,
		animalRequest.Ds,
		animalRequest.Diet,
		animalRequest.Description,
		animalRequest.Category,
	)

	if err != nil {
		w.Write(animalRequestErrorResponse(w, err))
		return
	}

	w.Write(animalRequestSuccessResponse(w, animalRequest))

}

func GetAnimals(w http.ResponseWriter, r *http.Request) {
	// Check if it's an OPTIONS request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
	w.Header().Set("Content-Type", "application/json")

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
		responses.CustomResponse(w, nil, e.Error(), http.StatusUnauthorized, e.Error())
		return
	}

	category := r.URL.Query().Get("category")
	pageStr := r.URL.Query().Get("page")

	id := r.URL.Query().Get("id")

	var response []byte
	if id != "" {
		response = getAnimalById(w, r, id)
	} else {
		response = getAnimalsByCategory(w, r, category, pageStr)
	}

	if response != nil {
		w.Write(response)
	}
}

func getAnimalsByCategory(w http.ResponseWriter, r *http.Request, category string, pageStr string) []byte {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Type", "application/json")

	if category == "" {
		resp, err := responses.MissingURLParametersResponse(w)
		if err != nil {
			return nil
		}
		return resp
	}

	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		resp, err := responses.CustomResponse(w, nil, "Invalid value for page parameter", http.StatusBadRequest, err.Error())
		if err != nil {
			return nil
		}
		return resp
	}

	// Set up pagination parameters (you can customize these)
	itemsPerPage := 10
	offset := (page - 1) * itemsPerPage

	db := r.Context().Value("db").(*sql.DB)

	// Query the database
	query := "SELECT id,photos,name FROM animals WHERE category = ? LIMIT ? OFFSET ?"
	rows, err := db.Query(query, category, itemsPerPage, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	defer rows.Close()

	// Fetch the entries
	var animalPreviews []AnimalPreview
	for rows.Next() {
		var entry AnimalPreview
		err := rows.Scan(&entry.ID, &entry.Photos, &entry.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}
		animalPreviews = append(animalPreviews, entry)
	}

	response, e := responses.CustomResponse(w, animalPreviews, pageStr, http.StatusOK, "")
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return nil
	}
	return response
}

func getAnimalById(w http.ResponseWriter, r *http.Request, id string) []byte {

	db := r.Context().Value("db").(*sql.DB)

	// Query the database
	query := "SELECT * FROM animals WHERE id = ?"
	rows, err := db.Query(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	defer rows.Close()

	// Fetch the entries
	var animal Animal
	for rows.Next() {
		err := rows.Scan(
			&animal.ID,
			&animal.Photos,
			&animal.Name,
			&animal.Taxonomy,
			&animal.Etymology,
			&animal.Iucn,
			&animal.Geo,
			&animal.Migration,
			&animal.Habitat,
			&animal.Dimensions,
			&animal.Ds,
			&animal.Diet,
			&animal.Description,
			&animal.Category,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}
	}

	response, e := responses.CustomResponse(w, animal, "Animal fetched from DB", http.StatusOK, "")
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return nil
	}
	return response
}
