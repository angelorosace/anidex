package responses

import (
	"encoding/json"
	"net/http"
)

type HttpResponse struct {
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

func CustomResponse(w http.ResponseWriter, data interface{}, message string, statusCode int, err string) ([]byte, error) {
	resp := HttpResponse{
		Data:    data,
		Error:   err,
		Message: message,
		Status:  statusCode,
	}
	jsonResponse, e := json.Marshal(resp)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return nil, e
	}
	return jsonResponse, nil
}
