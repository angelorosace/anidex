package handlers

import (
	"anidex_api/api/helpers"
	responses "anidex_api/http/responses"
	"database/sql"
	"encoding/json"
	"net/http"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type categorySuccessResponse struct {
	CategoryData []Category `json:"categoryData"`
	Error        string     `json:"error"`
	Message      string     `json:"message"`
	Status       int        `json:"status"`
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")

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

	res, err := db.Query("SELECT * from categories")
	if err != nil {
		w.Write(animalRequestErrorResponse(w, err))
		return
	}
	defer res.Close()

	// Create a slice to hold the results
	var categories []Category

	// Iterate through the rows and scan data into the slice
	for res.Next() {
		var cat Category
		err := res.Scan(&cat.ID, &cat.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, cat)
	}

	// Check for errors from iterating over rows
	if err := res.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := categorySuccessResponse{
		CategoryData: categories,
		Message:      "Categories successfully fetched",
		Status:       http.StatusOK,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}
