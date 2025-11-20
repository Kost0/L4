package handler

import (
	"encoding/json"
	"net/http"
	"sort"
)

type RequestPayload struct {
	Numbers []int `json:"numbers"`
}

func SortNums(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	sort.Ints(payload.Numbers)

	json.NewEncoder(w).Encode(payload.Numbers)
}
