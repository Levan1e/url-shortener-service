package http

import (
	"encoding/json"
	"net/http"

	"github.com/Levan1e/url-shortener-service/internal/domain"
	"github.com/Levan1e/url-shortener-service/pkg/logger"
)

func ParseReq[T any](r *http.Request) (*T, error) {
	defer r.Body.Close()
	var target T
	if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
		return &target, domain.InvalidEntry
	}
	return &target, nil
}

func BuildResponse[T any](w http.ResponseWriter, content T) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(content); err != nil {
		SetHttpError(w, err)
		return
	}
	//w.WriteHeader(http.StatusOK)
}

func SetHttpError(w http.ResponseWriter, e error) {
	err, isHandled := e.(*domain.HttpError)
	if !isHandled {
		err = domain.InternalServerError
	}
	logger.ErrorKV("http handle error", "err", e)
	http.Error(w, err.Error(), err.Code)
}
