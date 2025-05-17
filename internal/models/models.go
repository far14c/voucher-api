package models

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	ErrEmptyName         = errors.New("name cannot be empty")
	ErrEmptyCode         = errors.New("code cannot be empty")
	ErrInvalidPointsCost = errors.New("points cost must be positive")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrNegativePoints    = errors.New("points balance cannot be negative")
	ErrExpiredVoucher    = errors.New("voucher has expired")
	ErrNoItems           = errors.New("redemption must have at least one item")
	ErrInvalidStatus     = errors.New("invalid redemption status")
)

type Brand struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (b *Brand) Validate() error {
	return validateBrandInternal(*b)
}

type Voucher struct {
	ID          int       `json:"id"`
	BrandID     int       `json:"brand_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PointsCost  int       `json:"points_cost"`
	IsActive    bool      `json:"is_active"`
	ValidUntil  time.Time `json:"valid_until"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (v *Voucher) Validate() error {
	return validateVoucherInternal(*v)
}

type Customer struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	PointsBalance int       `json:"points_balance"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (c *Customer) Validate() error {
	return validateCustomerInternal(*c)
}

type Redemption struct {
	ID              int              `json:"id"`
	CustomerID      int              `json:"customer_id"`
	TotalPointsCost int              `json:"total_points_cost"`
	Status          string           `json:"status"`
	Items           []RedemptionItem `json:"items"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

func (r *Redemption) Validate() error {
	return validateRedemptionInternal(*r)
}

type RedemptionItem struct {
	ID           int       `json:"id"`
	RedemptionID int       `json:"redemption_id"`
	VoucherID    int       `json:"voucher_id"`
	PointsCost   int       `json:"points_cost"`
	CreatedAt    time.Time `json:"created_at"`
}

// Request/Response structures
type CreateBrandRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateVoucherRequest struct {
	BrandID     int       `json:"brand_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PointsCost  int       `json:"points_cost"`
	ValidUntil  time.Time `json:"valid_until"`
}

type RedemptionRequest struct {
	CustomerID int   `json:"customer_id"`
	VoucherIDs []int `json:"voucher_ids"`
}

// Validation functions
func validateBrandInternal(b Brand) error {
	if strings.TrimSpace(b.Name) == "" {
		return ErrEmptyName
	}
	return nil
}

func validateVoucherInternal(v Voucher) error {
	if strings.TrimSpace(v.Code) == "" {
		return ErrEmptyCode
	}
	if strings.TrimSpace(v.Name) == "" {
		return ErrEmptyName
	}
	if v.PointsCost <= 0 {
		return ErrInvalidPointsCost
	}
	if !v.ValidUntil.IsZero() && v.ValidUntil.Before(time.Now()) {
		return ErrExpiredVoucher
	}
	return nil
}

func validateCustomerInternal(c Customer) error {
	if strings.TrimSpace(c.Name) == "" {
		return ErrEmptyName
	}
	if !isValidEmail(c.Email) {
		return ErrInvalidEmail
	}
	if c.PointsBalance < 0 {
		return ErrNegativePoints
	}
	return nil
}

func validateRedemptionInternal(r Redemption) error {
	if len(r.Items) == 0 {
		return ErrNoItems
	}
	if !isValidStatus(r.Status) {
		return ErrInvalidStatus
	}
	return nil
}

// Helper functions
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":   true,
		"completed": true,
		"cancelled": true,
		"failed":    true,
	}
	return validStatuses[strings.ToLower(status)]
}
