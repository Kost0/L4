package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Kost0/L4/internal/models"
)

func GetStats(w http.ResponseWriter, r *http.Request) {
	statsResponse, err := json.Marshal(models.AllStats)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(statsResponse)
}
