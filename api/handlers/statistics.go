package handlers

import (
	"anidex_api/api/helpers"
	responses "anidex_api/http/responses"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type Stat struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
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

	// Parse and handle URL query parameters
	queryParams := r.URL.Query()

	// Extract specific query parameters
	table := queryParams.Get("table")
	groupBy := queryParams.Get("groupBy")

	if table == "" || groupBy == "" {
		resp, err := responses.MissingURLParametersResponse(w)
		if err != nil {
			return
		}
		w.Write(resp)
		return
	}

	//retrieve DB from context
	db := r.Context().Value("db").(*sql.DB)

	res, err := db.Query(fmt.Sprintf("SELECT %s,COUNT(*) from %s group by %s", groupBy, table, groupBy))
	if err != nil {
		resp, err := responses.MySqlError(w, err)
		if err != nil {
			return
		}
		w.Write(resp)
		return
	}
	defer res.Close()

	// Create a slice to hold the results
	stats := make(map[string]Stat)

	// Iterate through the rows and scan data into the slice
	for res.Next() {
		var category string
		var count int
		err := res.Scan(&category, &count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stat := Stat{Category: category, Count: count}
		stats[category] = stat
	}

	// Check for errors from iterating over rows
	if err := res.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := responses.HttpResponse{
		Data:    stats,
		Message: "Stats successfully computed",
		Status:  http.StatusOK,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}
