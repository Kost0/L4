package handler

import (
	"io/ioutil"
	"net/http"
	"sort"
)

type RequestPayload struct {
	Numbers []int `json:"numbers"`
}

func SortNums(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = payload.UnmarshalJSON(buf)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sort.Ints(payload.Numbers)

	bytes, err := payload.MarshalJSON()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
