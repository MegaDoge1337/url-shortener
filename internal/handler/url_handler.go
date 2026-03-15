package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/megadoge1337/url-shortener/internal/domain"
	"github.com/megadoge1337/url-shortener/internal/service"
)

type UrlHandler struct {
	s *service.UrlService
}

func NewUrlHandler(service *service.UrlService) *UrlHandler {
	return &UrlHandler{
		s: service,
	}
}

type CreateUrlRequest struct {
	URL string `json:"url"`
}

type UrlResponse struct {
	ID    int    `json:"id"`
	URL   string `json:"url"`
	Alias string `json:"alias"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func (h *UrlHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUrlRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request body", slog.Any("error", err))
		respondWithError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if req.URL == "" {
		respondWithError(w, http.StatusBadRequest, "validation failed", "url is required")
		return
	}

	url := domain.URL{
		URL:   req.URL,
		Alias: uuid.New().String(),
	}

	createUrl, err := h.s.Create(url)
	if err != nil {
		slog.Error("failed to create url", slog.Any("error", err))
		respondWithError(w, http.StatusInternalServerError, "failed to create url", err.Error())
		return
	}

	response := UrlResponse{
		ID:    createUrl.ID,
		URL:   createUrl.URL,
		Alias: createUrl.Alias,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (h *UrlHandler) RedirectByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id", "id must be an integer")
		return
	}

	url, err := h.s.GetById(id)
	if err != nil {
		slog.Error("failed to get url by id", slog.Any("error", err))
		respondWithError(w, http.StatusNotFound, "not found", "url not found")
		return
	}

	http.Redirect(w, r, normalizeURL(url.URL), http.StatusFound)
}

func (h *UrlHandler) RedirectByAlias(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")
	if alias == "" {
		respondWithError(w, http.StatusBadRequest, "invalid alias", "alias is required")
		return
	}

	url, err := h.s.GetByAlias(alias)
	if err != nil {
		slog.Error("failed to get url by alias", slog.Any("error", err))
		respondWithError(w, http.StatusNotFound, "not found", "url not found")
		return
	}

	http.Redirect(w, r, normalizeURL(url.URL), http.StatusFound)
}

func normalizeURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" {
		return "https://" + rawURL
	}
	return rawURL
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("failed to encode response", slog.Any("error", err))
	}
}

func respondWithError(w http.ResponseWriter, statusCode int, error string, message string) {
	response := ErrorResponse{
		Error:   error,
		Message: message,
	}
	respondWithJSON(w, statusCode, response)
}
