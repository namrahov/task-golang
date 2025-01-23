package util

import (
	"encoding/json"
	"net/http"
	"task-golang/model"
)

func ErrorRespondWriterJSON(w http.ResponseWriter, errRegister *model.ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errRegister.Code)
	json.NewEncoder(w).Encode(errRegister)
}
