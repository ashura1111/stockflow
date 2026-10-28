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

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(db *pgxpool.Pool) *ProductHandler {
	return &ProductHandler{
		service: services.NewProductService(db),
	}
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid product ID")
		return
	}

	product, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "Product not found")
		return
	}

	respondJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")
	lowStockStr := r.URL.Query().Get("low_stock")
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

	var lowStock *bool
	if lowStockStr == "true" {
		t := true
		lowStock = &t
	}

	filter := &models.ProductFilter{
		Search:   search,
		Status:   status,
		LowStock: lowStock,
		Page:     page,
		Limit:    limit,
	}

	list, err := h.service.ListProducts(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list products")
		return
	}

	respondJSON(w, http.StatusOK, list)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	product, err := h.service.CreateProduct(r.Context(), &req)
	if err != nil {
		if err.Error() == "SKU_ALREADY_EXISTS" {
			respondError(w, http.StatusConflict, "SKU_ALREADY_EXISTS", "SKU already exists")
			return
		}
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid product ID")
		return
	}

	var req models.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	product, err := h.service.UpdateProduct(r.Context(), id, &req)
	if err != nil {
		if err.Error() == "SKU_ALREADY_EXISTS" {
			respondError(w, http.StatusConflict, "SKU_ALREADY_EXISTS", "SKU already exists")
			return
		}
		if err.Error() == "product not found" {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Product not found")
			return
		}
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid product ID")
		return
	}

	err = h.service.DeleteProduct(r.Context(), id)
	if err != nil {
		if err.Error() == "PRODUCT_HAS_HISTORY" {
			respondError(w, http.StatusConflict, "PRODUCT_HAS_HISTORY", "Cannot delete product with history")
			return
		}
		if strings.Contains(err.Error(), "no rows") {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Product not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete product")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) AdjustStock(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid product ID")
		return
	}

	var req models.StockMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	product, err := h.service.AdjustStock(r.Context(), id, &req)
	if err != nil {
		if strings.HasPrefix(err.Error(), "INSUFFICIENT_STOCK") {
			parts := strings.Split(err.Error(), ":")
			available := 0
			if len(parts) > 1 {
				available, _ = strconv.Atoi(parts[1])
			}
			respondErrorWithAvailable(w, http.StatusConflict, "INSUFFICIENT_STOCK", "Insufficient stock", available)
			return
		}
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Product not found")
			return
		}
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) GetMovementHistory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid product ID")
		return
	}

	product, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "Product not found")
		return
	}

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

	filter := &models.StockMovementFilter{
		ProductID: product.ID,
		Page:      page,
		Limit:     limit,
	}

	movements, err := h.service.GetMovementHistory(r.Context(), product.ID, filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch movement history")
		return
	}

	respondJSON(w, http.StatusOK, movements)
}
func (h *ProductHandler) ListMovements(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 20

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}

	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	filter := &models.StockMovementFilter{
		Page:   page,
		Limit:  limit,
		Type:   r.URL.Query().Get("type"),
		Search: r.URL.Query().Get("search"),
	}

	if filter.Type != "" && filter.Type != "in" && filter.Type != "out" {
		respondError(w, http.StatusBadRequest, "INVALID_TYPE", "Invalid movement type")
		return
	}

	movements, err := h.service.GetMovementHistory(r.Context(), 0, filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch stock movements")
		return
	}

	respondJSON(w, http.StatusOK, movements)
}

func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, statusCode int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := models.ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	json.NewEncoder(w).Encode(resp)
}

func respondErrorWithAvailable(w http.ResponseWriter, statusCode int, code string, message string, available int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := models.ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	resp.Error.Available = available
	json.NewEncoder(w).Encode(resp)
}
