CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    productName VARCHAR(100) NOT NULL,
    description VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL
);