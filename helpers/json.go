package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func ReadFromRequestBody(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing body request")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteToResponseBody(w http.ResponseWriter, statusCode int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(value)
}

func WriteErrorResponse(w http.ResponseWriter, code int, status string, message string) {
	errorResponse := ErrorResponse{
		Code:    code,
		Message: message,
		Status:  status,
	}
	WriteToResponseBody(w, code, errorResponse)
}

func JSONFormatting(data any) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(jsonData))
}
