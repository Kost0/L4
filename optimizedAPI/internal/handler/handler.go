package handler

import (
	"encoding/json"
	"net/http"
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

	bubbleSort(payload.Numbers)

	json.NewEncoder(w).Encode(payload.Numbers)
}

func bubbleSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}
