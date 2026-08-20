CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    owner_id INT REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'confirmed' CHECK (status IN ('confirmed','processing','shipped','out_for_delivery','delivered','cancelled','return_requested','returned','refunded')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);

CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE RESTRICT,
    quantity INT NOT NULL CHECK (quantity > 0)
);


ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'customer' CHECK (role IN ('customer','seller','admin'));
ALTER TABLE users DROP COLUMN seller;