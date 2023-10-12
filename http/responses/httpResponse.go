package responses

type HttpResponse struct {
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}
