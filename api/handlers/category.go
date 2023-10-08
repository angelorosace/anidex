package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Type", "application/json")

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

	fmt.Println(res)

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

	fmt.Println(categories)

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
