package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"stockflow/backend/internal/models"
	"stockflow/backend/internal/services"
)

type SalesOrderHandler struct {
	service *services.SalesOrderService
}

func NewSalesOrderHandler(db *pgxpool.Pool) *SalesOrderHandler {
	return &SalesOrderHandler{
		service: services.NewSalesOrderService(db),
	}
}

func (h *SalesOrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid order ID")
		return
	}

	order, err := h.service.GetOrder(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "Order not found")
		return
	}

	respondJSON(w, http.StatusOK, order)
}

func (h *SalesOrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	filter := &models.OrderFilter{
		Search: search,
		Status: status,
		Page:   page,
		Limit:  limit,
	}

	list, err := h.service.ListOrders(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list orders")
		return
	}

	respondJSON(w, http.StatusOK, list)
}

func (h *SalesOrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	order, err := h.service.CreateOrder(r.Context(), &req)
	if err != nil {
		errMsg := err.Error()
		if strings.HasPrefix(errMsg, "INSUFFICIENT_STOCK") {
			parts := strings.Split(errMsg, ":")
			available := 0
			sku := ""
			if len(parts) > 1 {
				sku = parts[1]
			}
			if len(parts) > 2 {
				available, _ = strconv.Atoi(parts[2])
			}
			respondErrorWithSKUAndAvailable(w, http.StatusConflict, "INSUFFICIENT_STOCK", "Insufficient stock", sku, available)
			return
		}
		if strings.Contains(errMsg, "product not found") || strings.Contains(errMsg, "inactive product") {
			respondError(w, http.StatusBadRequest, "INVALID_ITEM", errMsg)
			return
		}
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", errMsg)
		return
	}

	respondJSON(w, http.StatusCreated, order)
}

func (h *SalesOrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid order ID")
		return
	}

	response, err := h.service.CancelOrder(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Order not found")
			return
		}
		if strings.Contains(err.Error(), "already cancelled") {
			respondError(w, http.StatusConflict, "ALREADY_CANCELLED", "Order is already cancelled")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *SalesOrderHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetDashboardStats(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch dashboard stats")
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

func respondErrorWithSKUAndAvailable(w http.ResponseWriter, statusCode int, code string, message string, sku string, available int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := models.ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	resp.Error.SKU = sku
	resp.Error.Available = available
	json.NewEncoder(w).Encode(resp)
}
