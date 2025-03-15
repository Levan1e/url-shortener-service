package v1

import (
	"context"
	"net/http"

	"github.com/Levan1e/url-shortener-service/internal/domain"
	"github.com/Levan1e/url-shortener-service/internal/models"
	http_helpers "github.com/Levan1e/url-shortener-service/pkg/http"
	"github.com/go-chi/chi/v5"
)

type ShortenerService interface {
	GetShortenByOriginal(context.Context, string) (string, error)
	GetOriginalByShorten(context.Context, string) (string, error)
}

type Handler struct {
	shortenerService ShortenerService
}

func NewHandler(shortener ShortenerService) *Handler {
	return &Handler{
		shortenerService: shortener,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/v1", func(r chi.Router) {
		r.Post("/shorten", h.CreateShortURL)
		r.Get("/{url}", h.GetOriginalURL)
	})
}

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	req, err := http_helpers.ParseReq[models.GetShortenByOriginalRequest](r)
	if err != nil {
		http_helpers.SetHttpError(w, err)
		return
	}
	shortURL, err := h.shortenerService.GetShortenByOriginal(r.Context(), req.Url)
	if err != nil {
		http_helpers.SetHttpError(w, err)
		return
	}
	http_helpers.BuildResponse(w, &models.GetShortenByOriginalResponse{ShortenUrl: shortURL})
}

func (h *Handler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "url")
	if shortURL == "" {
		http_helpers.SetHttpError(w, domain.InvalidEntry)
		return
	}
	original, err := h.shortenerService.GetOriginalByShorten(r.Context(), shortURL)
	if err != nil {
		http_helpers.SetHttpError(w, err)
		return
	}
	http_helpers.BuildResponse(w, &models.GetOriginalByShortenResponse{Url: original})
}
