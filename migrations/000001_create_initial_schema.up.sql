CREATE TABLE brands (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME
);

CREATE TABLE vouchers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    brand_id INT NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    points_cost INT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    valid_until DATETIME NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

CREATE TABLE customers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    points_balance INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME
);

CREATE TABLE redemptions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,
    total_points_cost INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);

CREATE TABLE redemption_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    redemption_id INT NOT NULL,
    voucher_id INT NOT NULL,
    points_cost INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (redemption_id) REFERENCES redemptions(id),
    FOREIGN KEY (voucher_id) REFERENCES vouchers(id)
);

-- Add indexes for better query performance
CREATE INDEX idx_vouchers_brand_id ON vouchers(brand_id);
CREATE INDEX idx_redemption_items_redemption_id ON redemption_items(redemption_id);
CREATE INDEX idx_redemption_items_voucher_id ON redemption_items(voucher_id); 