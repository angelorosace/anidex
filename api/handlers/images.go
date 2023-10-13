package handlers

import (
	"anidex_api/http/responses"
	"net/http"
)

func GetImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	//w.Header().Set("Content-Type", "application/json")
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
