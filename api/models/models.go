package models

type SeatListItem struct {
	Seat_class      string
	Seat_identifier string
	Is_booked       bool
	Price           string
}

type SeatsList struct {
	Seats []SeatListItem
}

type BookingPayload struct {
	SeatIDs     []string
	Name        string
	PhoneNumber string
}

type BookingReceipt struct {
	BookingID   int
	TotalAmount string
}

type SeatPrice struct {
	SeatClass   string
	MinPrice    []uint8
	NormalPrice []uint8
	MaxPrice    []uint8
}

type BookingDetails struct {
	BookingID int
	Seats     []SeatListItem
}

type BookingHistory struct {
	Bookings []BookingDetails
}
