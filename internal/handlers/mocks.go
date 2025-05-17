package handlers

import (
	"database/sql"
	"voucher-api/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockDB is a mock database implementation for testing
type MockDB struct {
	mock.Mock
}

// Mock database methods
func (m *MockDB) CreateBrand(brand *models.Brand) (int, error) {
	args := m.Called(brand)
	return args.Int(0), args.Error(1)
}

func (m *MockDB) GetBrand(id int) (*models.Brand, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Brand), args.Error(1)
}

func (m *MockDB) ListBrands() ([]models.Brand, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Brand), args.Error(1)
}

func (m *MockDB) CreateVoucher(voucher *models.Voucher) (int, error) {
	args := m.Called(voucher)
	return args.Int(0), args.Error(1)
}

func (m *MockDB) GetVoucher(id int) (*models.Voucher, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Voucher), args.Error(1)
}

func (m *MockDB) ListVouchers() ([]models.Voucher, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Voucher), args.Error(1)
}

func (m *MockDB) GetCustomer(id int) (*models.Customer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockDB) CreateRedemption(redemption *models.Redemption) (int, error) {
	args := m.Called(redemption)
	return args.Int(0), args.Error(1)
}

func (m *MockDB) GetRedemption(id int) (*models.Redemption, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Redemption), args.Error(1)
}

func (m *MockDB) UpdateCustomerPoints(customerID int, points int) error {
	args := m.Called(customerID, points)
	return args.Error(0)
}

func (m *MockDB) BeginTx() (*sql.Tx, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sql.Tx), args.Error(1)
}

func (m *MockDB) GetVouchersByBrand(brandID int) ([]models.Voucher, error) {
	args := m.Called(brandID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Voucher), args.Error(1)
}
