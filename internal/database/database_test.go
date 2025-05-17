package database

import (
	"database/sql"
	"testing"
	"time"
	"voucher-api/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetVouchersByBrand(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	// Create a new DB instance with the mock database
	dbInstance := NewDB(db)

	// Test cases
	tests := []struct {
		name      string
		brandID   int
		mockSetup func()
		want      []models.Voucher
		wantErr   bool
	}{
		{
			name:    "successful retrieval",
			brandID: 1,
			mockSetup: func() {
				now := time.Now()
				validUntil := now.Add(24 * time.Hour)
				rows := sqlmock.NewRows([]string{
					"id", "brand_id", "code", "name", "description",
					"points_cost", "is_active", "valid_until", "created_at", "updated_at",
				}).AddRow(
					1, 1, "CODE1", "Test Voucher 1", "Description 1",
					100, true, validUntil, now, now,
				).AddRow(
					2, 1, "CODE2", "Test Voucher 2", "Description 2",
					200, true, validUntil, now, now,
				)

				mock.ExpectQuery(`
					SELECT id, brand_id, code, name, description, points_cost, 
					       is_active, valid_until, created_at, updated_at
					FROM vouchers 
					WHERE brand_id = ?`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: []models.Voucher{
				{
					ID:          1,
					BrandID:     1,
					Code:        "CODE1",
					Name:        "Test Voucher 1",
					Description: "Description 1",
					PointsCost:  100,
					IsActive:    true,
				},
				{
					ID:          2,
					BrandID:     1,
					Code:        "CODE2",
					Name:        "Test Voucher 2",
					Description: "Description 2",
					PointsCost:  200,
					IsActive:    true,
				},
			},
			wantErr: false,
		},
		{
			name:    "no vouchers found",
			brandID: 2,
			mockSetup: func() {
				mock.ExpectQuery(`
					SELECT id, brand_id, code, name, description, points_cost, 
					       is_active, valid_until, created_at, updated_at
					FROM vouchers 
					WHERE brand_id = ?`).
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "brand_id", "code", "name", "description",
						"points_cost", "is_active", "valid_until", "created_at", "updated_at",
					}))
			},
			want:    []models.Voucher{},
			wantErr: false,
		},
		{
			name:    "database error",
			brandID: 3,
			mockSetup: func() {
				mock.ExpectQuery(`
					SELECT id, brand_id, code, name, description, points_cost, 
					       is_active, valid_until, created_at, updated_at
					FROM vouchers 
					WHERE brand_id = ?`).
					WithArgs(3).
					WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.mockSetup()

			// Call the function being tested
			got, err := dbInstance.GetVouchersByBrand(tt.brandID)

			// Check error expectations
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// For successful cases, compare the results
			// Note: We only compare specific fields since time fields will be different
			for i, voucher := range got {
				assert.Equal(t, tt.want[i].ID, voucher.ID)
				assert.Equal(t, tt.want[i].BrandID, voucher.BrandID)
				assert.Equal(t, tt.want[i].Code, voucher.Code)
				assert.Equal(t, tt.want[i].Name, voucher.Name)
				assert.Equal(t, tt.want[i].Description, voucher.Description)
				assert.Equal(t, tt.want[i].PointsCost, voucher.PointsCost)
				assert.Equal(t, tt.want[i].IsActive, voucher.IsActive)
			}

			// Verify that all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
