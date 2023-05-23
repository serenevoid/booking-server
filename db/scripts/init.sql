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

CREATE TABLE receipt (
    booking_id SERIAL PRIMARY KEY,
    phone VARCHAR(12) NOT NULL,
    price MONEY NOT NULL,
);

CREATE TABLE booking (
    id SERIAL PRIMARY KEY,
    seat_identifier VARCHAR(12) UNIQUE NOT NULL,
    booking_id INT NOT NULL,
    name VARCHAR(20) NOT NULL,
    phone VARCHAR(12) NOT NULL,
    price MONEY NOT NULL,
    FOREIGN KEY (seat_identifier) REFERENCES seats(seat_identifier)
    FOREIGN KEY (booking_id) REFERENCES receipt(booking_id)
);
