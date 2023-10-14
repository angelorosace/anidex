package handlers

import (
	"anidex_api/http/responses"
	"net/http"
)

func GetImageByPath(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

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
