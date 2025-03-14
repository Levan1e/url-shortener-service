package api

import (
	"encoding/json"
	"net/http"

	"github.com/Levan1e/url-shortener-service/internal/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	shortener *service.ShortenerService
}

func NewHandler(shortener *service.ShortenerService) *Handler {
	return &Handler{
		shortener: shortener,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/shorten", h.CreateShortURL)
	r.Get("/{shortURL}", h.GetOriginalURL)
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

type originalResponse struct {
	OriginalURL string `json:"original_url"`
}

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	shortURL, err := h.shortener.Shorten(req.URL)
	if err != nil {
		http.Error(w, "Ошибка генерации короткого URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := shortenResponse{ShortURL: shortURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	if shortURL == "" {
		http.Error(w, "Короткий URL не указан", http.StatusBadRequest)
		return
	}

	original, err := h.shortener.Resolve(shortURL)
	if err != nil || original == "" {
		http.Error(w, "Оригинальный URL не найден", http.StatusNotFound)
		return
	}

	resp := originalResponse{OriginalURL: original}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
