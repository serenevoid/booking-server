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

CREATE TABLE booking (
    booking_id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    phone VARCHAR(20) NOT NULL
);

CREATE TABLE receipt (
    booking_id INT REFERENCES booking(booking_id),
    seat_identifier VARCHAR(12) REFERENCES seats(seat_identifier),
    price MONEY NOT NULL,
    PRIMARY KEY (booking_id, seat_identifier)
);
