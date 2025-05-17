package models

import (
	"testing"
	"time"
)

func TestBrand_Validate(t *testing.T) {
	tests := []struct {
		name    string
		brand   Brand
		wantErr bool
	}{
		{
			name: "valid brand",
			brand: Brand{
				Name:        "Test Brand",
				Description: "Test Description",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			brand: Brand{
				Name:        "",
				Description: "Test Description",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateBrand(tt.brand); (err != nil) != tt.wantErr {
				t.Errorf("Brand.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVoucher_Validate(t *testing.T) {
	validTime := time.Now().Add(24 * time.Hour)
	invalidTime := time.Now().Add(-24 * time.Hour)

	tests := []struct {
		name    string
		voucher Voucher
		wantErr bool
	}{
		{
			name: "valid voucher",
			voucher: Voucher{
				BrandID:    1,
				Code:       "TEST123",
				Name:       "Test Voucher",
				PointsCost: 100,
				IsActive:   true,
				ValidUntil: validTime,
			},
			wantErr: false,
		},
		{
			name: "invalid code",
			voucher: Voucher{
				BrandID:    1,
				Code:       "",
				Name:       "Test Voucher",
				PointsCost: 100,
				ValidUntil: validTime,
			},
			wantErr: true,
		},
		{
			name: "invalid points cost",
			voucher: Voucher{
				BrandID:    1,
				Code:       "TEST123",
				Name:       "Test Voucher",
				PointsCost: -1,
				ValidUntil: validTime,
			},
			wantErr: true,
		},
		{
			name: "expired voucher",
			voucher: Voucher{
				BrandID:    1,
				Code:       "TEST123",
				Name:       "Test Voucher",
				PointsCost: 100,
				ValidUntil: invalidTime,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateVoucher(tt.voucher); (err != nil) != tt.wantErr {
				t.Errorf("Voucher.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCustomer_Validate(t *testing.T) {
	tests := []struct {
		name     string
		customer Customer
		wantErr  bool
	}{
		{
			name: "valid customer",
			customer: Customer{
				Name:          "John Doe",
				Email:         "john@example.com",
				PointsBalance: 0,
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			customer: Customer{
				Name:          "John Doe",
				Email:         "invalid-email",
				PointsBalance: 0,
			},
			wantErr: true,
		},
		{
			name: "negative points balance",
			customer: Customer{
				Name:          "John Doe",
				Email:         "john@example.com",
				PointsBalance: -100,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateCustomer(tt.customer); (err != nil) != tt.wantErr {
				t.Errorf("Customer.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedemption_Validate(t *testing.T) {
	tests := []struct {
		name       string
		redemption Redemption
		wantErr    bool
	}{
		{
			name: "valid redemption",
			redemption: Redemption{
				CustomerID:      1,
				TotalPointsCost: 100,
				Status:          "pending",
				Items: []RedemptionItem{
					{
						VoucherID:  1,
						PointsCost: 100,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "no items",
			redemption: Redemption{
				CustomerID:      1,
				TotalPointsCost: 100,
				Status:          "pending",
				Items:           []RedemptionItem{},
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			redemption: Redemption{
				CustomerID:      1,
				TotalPointsCost: 100,
				Status:          "invalid_status",
				Items: []RedemptionItem{
					{
						VoucherID:  1,
						PointsCost: 100,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRedemption(tt.redemption); (err != nil) != tt.wantErr {
				t.Errorf("Redemption.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper validation functions that would be implemented in models.go
func validateBrand(b Brand) error {
	// Add implementation
	return nil
}

func validateVoucher(v Voucher) error {
	// Add implementation
	return nil
}

func validateCustomer(c Customer) error {
	// Add implementation
	return nil
}

func validateRedemption(r Redemption) error {
	// Add implementation
	return nil
}
