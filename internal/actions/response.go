package actions

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data,omitempty"`
}

func (r *JSONResponse) PrepareResponse(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(r)
}

func (r *JSONResponse) String() string {
	body, _ := json.Marshal(r)
	return string(body)
}

func successData(data any) *JSONResponse {
	return &JSONResponse{Success: true, Data: data}
}
