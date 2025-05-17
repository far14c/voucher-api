package database

import (
	"database/sql"
	"voucher-api/internal/models"
)

// DB holds the database connection
type DB struct {
	db *sql.DB
}

// NewDB creates a new DB instance
func NewDB(db *sql.DB) *DB {
	return &DB{db: db}
}

// GetVouchersByBrand retrieves all vouchers for a given brand ID
func (d *DB) GetVouchersByBrand(brandID int) ([]models.Voucher, error) {
	query := `
		SELECT id, brand_id, code, name, description, points_cost, 
		       is_active, valid_until, created_at, updated_at
		FROM vouchers 
		WHERE brand_id = ?`

	rows, err := d.db.Query(query, brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vouchers []models.Voucher
	for rows.Next() {
		var v models.Voucher
		err := rows.Scan(
			&v.ID,
			&v.BrandID,
			&v.Code,
			&v.Name,
			&v.Description,
			&v.PointsCost,
			&v.IsActive,
			&v.ValidUntil,
			&v.CreatedAt,
			&v.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		vouchers = append(vouchers, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return vouchers, nil
}

// Close closes the database connection
func (d *DB) Close() error {
	return d.db.Close()
}

// BeginTx starts a new transaction
func (d *DB) BeginTx() (*sql.Tx, error) {
	return d.db.Begin()
}

// CreateBrand creates a new brand
func (d *DB) CreateBrand(brand *models.Brand) (int, error) {
	query := `INSERT INTO brands (name, description) VALUES (?, ?)`
	result, err := d.db.Exec(query, brand.Name, brand.Description)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// GetBrand retrieves a brand by ID
func (d *DB) GetBrand(id int) (*models.Brand, error) {
	var brand models.Brand
	err := d.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM brands WHERE id = ?", id).
		Scan(&brand.ID, &brand.Name, &brand.Description, &brand.CreatedAt, &brand.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &brand, nil
}

// ListBrands retrieves all brands
func (d *DB) ListBrands() ([]models.Brand, error) {
	rows, err := d.db.Query("SELECT id, name, description, created_at, updated_at FROM brands")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brands []models.Brand
	for rows.Next() {
		var b models.Brand
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		brands = append(brands, b)
	}
	return brands, rows.Err()
}

// CreateVoucher creates a new voucher
func (d *DB) CreateVoucher(voucher *models.Voucher) (int, error) {
	query := `INSERT INTO vouchers (brand_id, code, name, description, points_cost, is_active, valid_until) 
	         VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := d.db.Exec(query, voucher.BrandID, voucher.Code, voucher.Name,
		voucher.Description, voucher.PointsCost, voucher.IsActive, voucher.ValidUntil)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// GetVoucher retrieves a voucher by ID
func (d *DB) GetVoucher(id int) (*models.Voucher, error) {
	var v models.Voucher
	err := d.db.QueryRow(`SELECT id, brand_id, code, name, description, points_cost, 
		is_active, valid_until, created_at, updated_at FROM vouchers WHERE id = ?`, id).
		Scan(&v.ID, &v.BrandID, &v.Code, &v.Name, &v.Description, &v.PointsCost,
			&v.IsActive, &v.ValidUntil, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// ListVouchers retrieves all vouchers
func (d *DB) ListVouchers() ([]models.Voucher, error) {
	rows, err := d.db.Query(`SELECT id, brand_id, code, name, description, points_cost, 
		is_active, valid_until, created_at, updated_at FROM vouchers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vouchers []models.Voucher
	for rows.Next() {
		var v models.Voucher
		if err := rows.Scan(&v.ID, &v.BrandID, &v.Code, &v.Name, &v.Description, &v.PointsCost,
			&v.IsActive, &v.ValidUntil, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		vouchers = append(vouchers, v)
	}
	return vouchers, rows.Err()
}

// GetCustomer retrieves a customer by ID
func (d *DB) GetCustomer(id int) (*models.Customer, error) {
	var c models.Customer
	err := d.db.QueryRow("SELECT id, name, email, points_balance, created_at, updated_at FROM customers WHERE id = ?", id).
		Scan(&c.ID, &c.Name, &c.Email, &c.PointsBalance, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// CreateRedemption creates a new redemption
func (d *DB) CreateRedemption(redemption *models.Redemption) (int, error) {
	query := `INSERT INTO redemptions (customer_id, total_points_cost, status) VALUES (?, ?, ?)`
	result, err := d.db.Exec(query, redemption.CustomerID, redemption.TotalPointsCost, redemption.Status)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// GetRedemption retrieves a redemption by ID
func (d *DB) GetRedemption(id int) (*models.Redemption, error) {
	var r models.Redemption
	err := d.db.QueryRow(`SELECT id, customer_id, total_points_cost, status, created_at, updated_at 
		FROM redemptions WHERE id = ?`, id).
		Scan(&r.ID, &r.CustomerID, &r.TotalPointsCost, &r.Status, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// UpdateCustomerPoints updates a customer's points balance
func (d *DB) UpdateCustomerPoints(customerID int, points int) error {
	_, err := d.db.Exec("UPDATE customers SET points_balance = ? WHERE id = ?", points, customerID)
	return err
}
