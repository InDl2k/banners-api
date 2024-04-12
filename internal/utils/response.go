package utils

import (
	"banners/internal/database/models"
	"encoding/json"
	"net/http"
)

func MessageError(error string) *map[string]interface{} {
	return &map[string]interface{}{"error": error}
}

func Respond(w http.ResponseWriter, data *map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&data)
}

func RespondBanners(w http.ResponseWriter, array *[]models.Banner) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&array)
}
