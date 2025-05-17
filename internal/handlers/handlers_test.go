package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"voucher-api/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBrand(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		setupMock      func(*MockDB)
	}{
		{
			name: "valid brand creation",
			requestBody: map[string]interface{}{
				"name":        "Test Brand",
				"description": "Test Description",
			},
			expectedStatus: http.StatusCreated,
			setupMock: func(m *MockDB) {
				m.On("CreateBrand", mock.Anything).Return(1, nil)
			},
		},
		{
			name: "invalid brand - empty name",
			requestBody: map[string]interface{}{
				"name":        "",
				"description": "Test Description",
			},
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockDB) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB)

			handler := NewHandler(mockDB)
			router := chi.NewRouter()
			router.Post("/brands", handler.CreateBrand)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/brands", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestGetBrand(t *testing.T) {
	tests := []struct {
		name           string
		brandID        string
		expectedStatus int
		setupMock      func(*MockDB)
	}{
		{
			name:           "existing brand",
			brandID:        "1",
			expectedStatus: http.StatusOK,
			setupMock: func(m *MockDB) {
				m.On("GetBrand", 1).Return(&models.Brand{
					ID:          1,
					Name:        "Test Brand",
					Description: "Test Description",
				}, nil)
			},
		},
		{
			name:           "non-existent brand",
			brandID:        "999",
			expectedStatus: http.StatusNotFound,
			setupMock: func(m *MockDB) {
				m.On("GetBrand", 999).Return(nil, sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB)

			handler := NewHandler(mockDB)
			router := chi.NewRouter()
			router.Get("/brands/{id}", handler.GetBrand)

			req := httptest.NewRequest("GET", "/brands/"+tt.brandID, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestCreateVoucher(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		setupMock      func(*MockDB)
	}{
		{
			name: "valid voucher creation",
			requestBody: map[string]interface{}{
				"brand_id":    1,
				"code":        "TEST123",
				"name":        "Test Voucher",
				"description": "Test Description",
				"points_cost": 100,
			},
			expectedStatus: http.StatusCreated,
			setupMock: func(m *MockDB) {
				m.On("CreateVoucher", mock.Anything).Return(1, nil)
			},
		},
		{
			name: "invalid voucher - negative points",
			requestBody: map[string]interface{}{
				"brand_id":    1,
				"code":        "TEST123",
				"name":        "Test Voucher",
				"description": "Test Description",
				"points_cost": -100,
			},
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockDB) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB)

			handler := NewHandler(mockDB)
			router := chi.NewRouter()
			router.Post("/vouchers", handler.CreateVoucher)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/vouchers", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestCreateRedemption(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		setupMock      func(*MockDB)
	}{
		{
			name: "valid redemption",
			requestBody: map[string]interface{}{
				"customer_id": 1,
				"voucher_ids": []int{1, 2},
			},
			expectedStatus: http.StatusCreated,
			setupMock: func(m *MockDB) {
				m.On("GetCustomer", 1).Return(&models.Customer{ID: 1, PointsBalance: 1000}, nil)
				m.On("GetVoucher", 1).Return(&models.Voucher{ID: 1, PointsCost: 100}, nil)
				m.On("GetVoucher", 2).Return(&models.Voucher{ID: 2, PointsCost: 200}, nil)
				m.On("CreateRedemption", mock.Anything).Return(1, nil)
				m.On("UpdateCustomerPoints", 1, 700).Return(nil)
			},
		},
		{
			name: "insufficient points",
			requestBody: map[string]interface{}{
				"customer_id": 1,
				"voucher_ids": []int{1},
			},
			expectedStatus: http.StatusBadRequest,
			setupMock: func(m *MockDB) {
				m.On("GetCustomer", 1).Return(&models.Customer{ID: 1, PointsBalance: 50}, nil)
				m.On("GetVoucher", 1).Return(&models.Voucher{ID: 1, PointsCost: 100}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB)

			handler := NewHandler(mockDB)
			router := chi.NewRouter()
			router.Post("/redemptions", handler.CreateRedemption)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/redemptions", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockDB.AssertExpectations(t)
		})
	}
}

func TestGetVouchersByBrand(t *testing.T) {
	mockDB := new(MockDB)
	handler := NewHandler(mockDB)

	tests := []struct {
		name       string
		brandID    string
		setupMock  func()
		wantStatus int
		wantBody   string
	}{
		{
			name:    "success",
			brandID: "1",
			setupMock: func() {
				vouchers := []models.Voucher{
					{
						ID:          1,
						BrandID:     1,
						Code:        "CODE1",
						Name:        "Test Voucher 1",
						Description: "Test Description 1",
						PointsCost:  100,
						IsActive:    true,
					},
					{
						ID:          2,
						BrandID:     1,
						Code:        "CODE2",
						Name:        "Test Voucher 2",
						Description: "Test Description 2",
						PointsCost:  200,
						IsActive:    true,
					},
				}
				mockDB.On("GetVouchersByBrand", 1).Return(vouchers, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `[{"id":1,"brand_id":1,"code":"CODE1","name":"Test Voucher 1","description":"Test Description 1","points_cost":100,"is_active":true},{"id":2,"brand_id":1,"code":"CODE2","name":"Test Voucher 2","description":"Test Description 2","points_cost":200,"is_active":true}]`,
		},
		{
			name:    "invalid brand ID",
			brandID: "invalid",
			setupMock: func() {
				// No mock setup needed for invalid ID
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid brand ID\n",
		},
		{
			name:    "database error",
			brandID: "1",
			setupMock: func() {
				mockDB.On("GetVouchersByBrand", 1).Return(nil, sql.ErrNoRows)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "sql: no rows in result set\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest("GET", "/vouchers/brand?id="+tt.brandID, nil)
			w := httptest.NewRecorder()

			handler.GetVouchersByBrand(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
			assert.JSONEq(t, tt.wantBody, string(body))
		})
	}
}
