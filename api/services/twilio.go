package services

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sfreiberg/gotwilio"
)

func SendNotification(bookingID string, amount string, phone string) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	AccountSID, exists := os.LookupEnv("AccountSID")
	if !exists {
		return fmt.Errorf("cannot find AccountSID")
	}
	AuthToken, exists := os.LookupEnv("AuthToken")
	if !exists {
		return fmt.Errorf("cannot find AuthToken")
	}
	PhoneNumber, exists := os.LookupEnv("PhoneNumber")
	if !exists {
		return fmt.Errorf("cannot find PhoneNumber")
	}
	message := fmt.Sprintf("Thanks for using our service. Your booking ID is %s and the total amount is %s", bookingID, amount)
	twilio := gotwilio.NewTwilioClient(AccountSID, AuthToken)

	_, _, err = twilio.SendSMS(PhoneNumber, phone, message, "", "")

	return err
}
