package handlers

import (
	"anidex_api/api/helpers"
	"anidex_api/http/responses"
	"net/http"
)

func GetImageByPath(w http.ResponseWriter, r *http.Request) {
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
		responses.CustomResponse(w, nil, e.Error(), http.StatusUnauthorized, e.Error())
		return
	}

	// Extract the file name from the URL path
	photoPath := r.URL.Query().Get("photo")

	if photoPath == "" {
		resp, err := responses.MissingURLParametersResponse(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(resp)
		return
	}

	http.ServeFile(w, r, photoPath)
}
