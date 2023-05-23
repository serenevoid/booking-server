package main

import (
	"fmt"
	"log"

    "booking-server/api/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
    fmt.Println("=== Starting Booking API ===")
    app := fiber.New()

    app.Get("/", handlers.Welcome)
    app.Get("/seats", handlers.GetAllSeats)
    app.Get("/seats/:id", handlers.GetSeatById)
    app.Post("/booking", handlers.BookSeats)
    app.Get("/bookings", handlers.ViewBookings)

    log.Fatal(app.Listen(":8080"))
}
