package responses

import (
	"encoding/json"
	"net/http"
)

func MySqlError(w http.ResponseWriter, err error) ([]byte, error) {
	resp := HttpResponse{
		Error:   err.Error(),
		Message: "MySQL query not ran",
		Status:  http.StatusInternalServerError,
	}
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	return jsonResponse, nil
}
