package handlers

import (
	"booking-server/api/models"
	"booking-server/api/utils"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Welcome(c *fiber.Ctx) error {
	c.SendString("Welcome to Flurn Seat Booking Service!")
	return nil
}

func GetAllSeats(c *fiber.Ctx) error {
	list := utils.GetAllSeats()

	data, err := json.Marshal(list)
	if err != nil {
		log.Fatal(err)
	}
	seats_data := strings.ReplaceAll(string(data), ",\"Price\":\"\"", "")
	c.SendString(seats_data)
	return nil
}

func GetSeatById(c *fiber.Ctx) error {
	id := c.AllParams()["id"]
	if utils.SeatExists(id) {
		var seat_details models.SeatListItem
		seat_details.Is_booked = utils.IsSeatBooked(id)
		seat_details.Seat_identifier = id
		seat_details.Seat_class = utils.GetSeatClass(id)
		seat_details.Price = utils.GetPrice(seat_details.Seat_class)

		data, err := json.Marshal(seat_details)
		if err != nil {
			log.Fatal(err)
		}

		c.SendString(string(data))
	} else {
		c.SendStatus(404)
		c.SendString(fmt.Sprintf("Error: Seat %v does not exist", id))
	}
	return nil
}

func BookSeats(c *fiber.Ctx) error {
	var payload models.BookingPayload
	err := c.BodyParser(&payload)
	if err != nil {
		log.Panic(err)
	}

	// Check if any seats are duplicates
	seatMap := make(map[string]bool)
	for _, seat := range payload.SeatIDs {
		if seatMap[seat] {
			c.SendStatus(403)
			c.SendString(fmt.Sprintf("Error: Seat %v has duplicates in the payload", seat))
			return nil
		}
		seatMap[seat] = true
	}

	// Safety checks
	for _, seat_id := range payload.SeatIDs {
		if !utils.SeatExists(seat_id) {
			c.SendStatus(404)
			c.SendString(fmt.Sprintf("Error: Seat %v does not exist", seat_id))
			return nil
		}

		if utils.IsSeatBooked(seat_id) {
			c.SendStatus(409)
			c.SendString(fmt.Sprintf("Error: Seat %v is already booked", seat_id))
			return nil
		}
	}

	receipt := utils.BookSeats(payload)

	data, err := json.Marshal(receipt)
	if err != nil {
		log.Fatal(err)
	}

	c.SendString(string(data))
	return nil
}

func ViewBookings(c *fiber.Ctx) error {
    phone := c.Query("userIdentifier")
    bookingHistory := utils.GetAllBookings(phone)

	data, err := json.Marshal(bookingHistory)
	if err != nil {
		log.Fatal(err)
	}

    filtered_data := strings.ReplaceAll(string(data), "\"Is_booked\":false,", "")

	c.SendString(filtered_data)
	return nil
}
