CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(100) NOT NULL,
    descripction VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL
);