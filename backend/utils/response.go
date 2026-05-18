package utils

import (
	"context"
	"encoding/json"
	"net/http"
)

func WriteResponse(ctx context.Context, w http.ResponseWriter, response interface{}, httpStatus int) {
	var resp []byte
	var err error
	if response != nil {
		resp, err = json.Marshal(response)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	if response != nil {
		_, err = w.Write(resp)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func Ptr[T any](v T) *T {
	return &v
}

func PtrOrZero(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func PtrOrZero64(f *float64) float64 {
	if f != nil {
		return *f
	}
	return 0
}
