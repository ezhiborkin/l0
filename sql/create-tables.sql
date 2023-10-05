-- Создаем таблицу "orders"
CREATE TABLE orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(255),
    locale VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey VARCHAR(10),
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard VARCHAR(10)
);

-- Создаем таблицу "delivery"
CREATE TABLE delivery (
    order_uid VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    phone VARCHAR(255),
    zip VARCHAR(255),
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255)
);

-- Создаем таблицу "payment"
CREATE TABLE payment (
    order_uid VARCHAR(255) PRIMARY KEY,
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(255),
    provider VARCHAR(255),
    amount INT,
    payment_dt BIGINT,
    bank VARCHAR(255),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

-- Создаем таблицу "items"
CREATE TABLE items (
    order_uid VARCHAR(255) PRIMARY KEY,
    chrt_id INT,
    track_number VARCHAR(255),
    price INT,
    rid VARCHAR(255),
    name VARCHAR(255),
    sale INT,
    size VARCHAR(255),
    total_price INT,
    nm_id INT,
    brand VARCHAR(255),
    status INT
);

-- CREATE TABLE roflans (
--     order_uid VARCHAR(255) PRIMARY KEY,
--     track_number VARCHAR(255)
-- );
