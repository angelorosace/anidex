package responses

import (
	"encoding/json"
	"net/http"
)

func MissingURLParametersResponse(w http.ResponseWriter) ([]byte, error) {
	resp := HttpResponse{
		Error:   "Missing one or more Url Parameters",
		Message: "Unable to contact endpoint",
		Status:  http.StatusBadRequest,
	}
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	return jsonResponse, nil
}
