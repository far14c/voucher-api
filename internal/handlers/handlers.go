package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"voucher-api/internal/models"

	"github.com/go-chi/chi/v5"
)

// Database interface defines all database operations
type Database interface {
	CreateBrand(brand *models.Brand) (int, error)
	GetBrand(id int) (*models.Brand, error)
	ListBrands() ([]models.Brand, error)
	CreateVoucher(voucher *models.Voucher) (int, error)
	GetVoucher(id int) (*models.Voucher, error)
	ListVouchers() ([]models.Voucher, error)
	GetCustomer(id int) (*models.Customer, error)
	CreateRedemption(redemption *models.Redemption) (int, error)
	GetRedemption(id int) (*models.Redemption, error)
	UpdateCustomerPoints(customerID int, points int) error
	BeginTx() (*sql.Tx, error)
	GetVouchersByBrand(brandID int) ([]models.Voucher, error)
}

// Handler holds the HTTP handlers and db connection
type Handler struct {
	db Database
}

// NewHandler creates a new handler with the given database
func NewHandler(db Database) *Handler {
	return &Handler{db: db}
}

// CreateBrand handles brand creation
func (h *Handler) CreateBrand(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBrandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	brand := &models.Brand{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := brand.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.db.CreateBrand(brand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// GetBrand handles retrieving a brand by ID
func (h *Handler) GetBrand(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid brand ID", http.StatusBadRequest)
		return
	}

	brand, err := h.db.GetBrand(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(brand)
}

// ListBrands handles retrieving all brands
func (h *Handler) ListBrands(w http.ResponseWriter, r *http.Request) {
	brands, err := h.db.ListBrands()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(brands)
}

// CreateVoucher handles voucher creation
func (h *Handler) CreateVoucher(w http.ResponseWriter, r *http.Request) {
	var req models.CreateVoucherRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	voucher := &models.Voucher{
		BrandID:     req.BrandID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		PointsCost:  req.PointsCost,
		ValidUntil:  req.ValidUntil,
		IsActive:    true,
	}

	if err := voucher.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.db.CreateVoucher(voucher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// GetVoucher handles retrieving a voucher by ID
func (h *Handler) GetVoucher(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid voucher ID", http.StatusBadRequest)
		return
	}

	voucher, err := h.db.GetVoucher(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(voucher)
}

// ListVouchers handles retrieving all vouchers
func (h *Handler) ListVouchers(w http.ResponseWriter, r *http.Request) {
	vouchers, err := h.db.ListVouchers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(vouchers)
}

// CreateRedemption handles voucher redemption
func (h *Handler) CreateRedemption(w http.ResponseWriter, r *http.Request) {
	var req models.RedemptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get customer
	customer, err := h.db.GetCustomer(req.CustomerID)
	if err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	// Calculate total points cost and validate vouchers
	var totalPoints int
	var items []models.RedemptionItem
	for _, vID := range req.VoucherIDs {
		voucher, err := h.db.GetVoucher(vID)
		if err != nil {
			http.Error(w, "Voucher not found", http.StatusNotFound)
			return
		}
		if !voucher.IsActive {
			http.Error(w, "Voucher is not active", http.StatusBadRequest)
			return
		}
		totalPoints += voucher.PointsCost
		items = append(items, models.RedemptionItem{
			VoucherID:  vID,
			PointsCost: voucher.PointsCost,
		})
	}

	if customer.PointsBalance < totalPoints {
		http.Error(w, "Insufficient points", http.StatusBadRequest)
		return
	}

	redemption := &models.Redemption{
		CustomerID:      req.CustomerID,
		TotalPointsCost: totalPoints,
		Status:          "pending",
		Items:           items,
	}

	id, err := h.db.CreateRedemption(redemption)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update customer points
	err = h.db.UpdateCustomerPoints(customer.ID, customer.PointsBalance-totalPoints)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// GetRedemption handles retrieving a redemption by ID
func (h *Handler) GetRedemption(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid redemption ID", http.StatusBadRequest)
		return
	}

	redemption, err := h.db.GetRedemption(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(redemption)
}

func (h *Handler) GetVouchersByBrand(w http.ResponseWriter, r *http.Request) {
	brandID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid brand ID", http.StatusBadRequest)
		return
	}

	vouchers, err := h.db.GetVouchersByBrand(brandID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(vouchers)
}
