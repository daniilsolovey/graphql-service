-- migrate:up
CREATE TABLE users (
  id SERIAL,
  name VARCHAR(50),
  phone VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE products (
  id SERIAL,
  name VARCHAR(50)
);

CREATE TABLE sms_codes (
  id SERIAL,
  phone VARCHAR(50) UNIQUE NOT NULL,
  code VARCHAR(50),
  expired_at TIMESTAMP
);

-- migrate:down
