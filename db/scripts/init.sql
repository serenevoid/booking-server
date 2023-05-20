CREATE TABLE pricing (
    id SERIAL PRIMARY KEY,
    seat_class VARCHAR(1) UNIQUE NOT NULL,
    min_price MONEY,
    normal_price MONEY,
    max_price MONEY
);

CREATE TABLE seats (
    id SERIAL PRIMARY KEY,
    seat_identifier VARCHAR(12) UNIQUE NOT NULL,
    seat_class VARCHAR(1) NOT NULL,
    FOREIGN KEY (seat_class) REFERENCES pricing(seat_class)
);
