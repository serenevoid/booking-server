package utils

import (
	"booking-server/api/models"
	"booking-server/api/services"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	HOST              = "db"
	PORT              = "5432"
	POSTGRES_USER     = "user"
	POSTGRES_PASSWORD = "pass"
	POSTGRES_DB       = "booking"
)

var DB *sql.DB

func init() {
	ConnectDB()
}

func ConnectDB() {
	connectionString := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=disable", POSTGRES_USER, POSTGRES_PASSWORD, HOST, PORT, POSTGRES_DB)
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("=== Connected to PostgresDB ===")
}

func CloseDB() {
	DB.Close()
}

func GetAllSeats() models.SeatsList {
	rows, err := DB.Query(`SELECT s.seat_class, s.seat_identifier, 
    CASE WHEN r.seat_identifier IS NULL THEN false ELSE true END AS is_booked FROM seats s 
    LEFT JOIN receipt r ON s.seat_identifier = r.seat_identifier ORDER BY s.seat_class;`)
	if err != nil {
		log.Panic(err)
	}

	var list models.SeatsList
	defer rows.Close()

	for rows.Next() {
		var seat models.SeatListItem
		err = rows.Scan(&seat.Seat_class, &seat.Seat_identifier, &seat.Is_booked)
		if err != nil {
			log.Panic(err)
		}
		list.Seats = append(list.Seats, seat)
	}

	return list
}

func SeatExists(seat_id string) bool {
	var is_booked int

	query := fmt.Sprintf("SELECT COUNT(*) FROM seats WHERE seat_identifier = '%s';", seat_id)
	err := DB.QueryRow(query).Scan(&is_booked)
	if err != nil {
		log.Panic(err)
	}

	return (is_booked == 1)
}

func IsSeatBooked(seat_id string) bool {
	var is_booked int

	query := fmt.Sprintf("SELECT COUNT(*) FROM receipt WHERE seat_identifier = '%s';", seat_id)
	err := DB.QueryRow(query).Scan(&is_booked)
	if err != nil {
		log.Panic(err)
	}

	return (is_booked == 1)
}

func GetSeatClass(seat_id string) string {
	var class string

	query := fmt.Sprintf("SELECT seat_class FROM seats WHERE seat_identifier = '%s';", seat_id)
	err := DB.QueryRow(query).Scan(&class)
	if err != nil {
		log.Panic(err)
	}

	return class
}

func GetPrice(class string) string {
	var total int
	var booked int
	Price := models.SeatPrice{
		SeatClass: class,
	}

	err := DB.QueryRow(fmt.Sprintf(`SELECT COALESCE(min_price, CAST(0 AS money)) AS min_price, 
    COALESCE(normal_price, CAST(0 AS money)) AS normal_price,
    COALESCE(max_price, CAST(0 AS money)) AS max_price 
    FROM pricing WHERE seat_class = '%s';`, class)).
		Scan(&Price.MinPrice, &Price.NormalPrice, &Price.MaxPrice)
	if err != nil {
		log.Panic(err)
	}

	err = DB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM seats WHERE seat_class = '%s';", class)).Scan(&total)
	if err != nil {
		log.Panic(err)
	}

	err = DB.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM receipt JOIN seats 
    ON seats.seat_identifier = receipt.seat_identifier WHERE seats.seat_class = '%s';`, class)).Scan(&booked)
	if err != nil {
		log.Panic(err)
	}

	fillRatio := 10 * booked / total

	if fillRatio < 4 {
		if string(Price.MinPrice) != "$0.00" {
			return string(Price.MinPrice)
		} else {
			return string(Price.NormalPrice)
		}
	} else if fillRatio <= 6 {
		return string(Price.NormalPrice)
	} else {
		if string(Price.MaxPrice) != "$0.00" {
			return string(Price.MaxPrice)
		} else {
			return string(Price.NormalPrice)
		}
	}
}

func BookSeats(payload models.BookingPayload) models.BookingReceipt {
	var id int
	query := fmt.Sprintf("INSERT INTO booking (booking_id, phone, name) VALUES (DEFAULT, '%s', '%s') RETURNING booking_id", payload.PhoneNumber, payload.Name)
	err := DB.QueryRow(query).Scan(&id)
	if err != nil {
		log.Panic(err)
	}

	for _, seat := range payload.SeatIDs {
		class := GetSeatClass(seat)
		price := GetPrice(class)
		query := fmt.Sprintf(`INSERT INTO receipt (booking_id, seat_identifier, price)
        VALUES (%d, '%s', '%s'::money)`, id, seat, price)
		_, err := DB.Exec(query)
		if err != nil {
			log.Panic(err)
		}
	}

	var price []uint8
	query_price := fmt.Sprintf(`SELECT SUM(price) FROM receipt 
        WHERE booking_id = '%d';`, id)
	err = DB.QueryRow(query_price).Scan(&price)
	if err != nil {
		log.Panic(err)
	}

	bill := models.BookingReceipt{
		BookingID:   id,
		TotalAmount: string(price),
	}

    err = services.SendNotification(fmt.Sprint(id), string(price), payload.PhoneNumber)
	if err != nil {
		log.Panic(err)
	}

	return bill
}

func GetAllBookings(phone string) models.BookingHistory {
	var data models.BookingHistory
	var bookingIDs []int
	rows, err := DB.Query(fmt.Sprintf("SELECT booking_id FROM booking WHERE phone = '%s';", phone))
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var bookingID int
		err = rows.Scan(&bookingID)
		if err != nil {
			log.Panic(err)
		}
		bookingIDs = append(bookingIDs, bookingID)
	}

	for _, bookingID := range bookingIDs {
		var bill models.BookingDetails
		query := fmt.Sprintf(`SELECT receipt.seat_identifier, receipt.price, seats.seat_class 
        FROM receipt JOIN seats ON receipt.seat_identifier = seats.seat_identifier 
        WHERE booking_id = %d ORDER BY receipt.booking_id;`, bookingID)
		rows, err := DB.Query(query)
		if err != nil {
			log.Panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var seat models.SeatListItem
			err = rows.Scan(&seat.Seat_identifier, &seat.Price, &seat.Seat_class)
			if err != nil {
				log.Panic(err)
			}
			bill.Seats = append(bill.Seats, seat)
		}
        bill.BookingID = bookingID
		data.Bookings = append(data.Bookings, bill)
	}

	return data
}
