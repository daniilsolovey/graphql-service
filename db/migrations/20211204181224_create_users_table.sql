-- migrate:up
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50),
  phone VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50)
);

CREATE TABLE sms_codes (
  id SERIAL PRIMARY KEY,
  phone VARCHAR(50) UNIQUE NOT NULL,
  code VARCHAR(50),
  expired_at TIMESTAMP
);

-- migrate:down
