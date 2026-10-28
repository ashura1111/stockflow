package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	if err := h.db.Ping(r.Context()); err != nil {
		respondError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Database connection failed")
		return
	}

	resp := map[string]string{
		"status": "ok",
	}
	respondJSON(w, http.StatusOK, resp)
}
