package handler

import (
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/valyala/fasthttp"
)

type RequestPayload struct {
	Numbers []int `json:"numbers"`
}

func RequestHandler(ctx *fasthttp.RequestCtx) {
	method := string(ctx.Method())
	path := string(ctx.Path())

	if path == "/sort" && method == "POST" {
		var payload RequestPayload

		buf := ctx.Request.Body()

		err := payload.UnmarshalJSON(buf)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		sort.Ints(payload.Numbers)

		bytes, err := payload.MarshalJSON()
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody(bytes)
	}
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
