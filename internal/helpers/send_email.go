package helpers

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendEmail sends an email with the OTP code to the specified email address.
// Requires EMAIL_FROM and EMAIL_PASSWORD environment variables.
func SendEmail(to, otp string) error {
	from := os.Getenv("EMAIL_FROM")
	if from == "" {
		return fmt.Errorf("EMAIL_FROM environment variable is not set")
	}

	password := os.Getenv("EMAIL_PASSWORD")
	if password == "" {
		return fmt.Errorf("EMAIL_PASSWORD environment variable is not set")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Compose the email message
	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP (One-Time Password) is: %s\n\nThis code is valid for 10 minutes. Do not share this code with anyone.", otp)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)

	// Create SMTP auth
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
