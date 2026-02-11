package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SendExistingOTP sends a provided OTP to the given phone number
func SendOTPToPhoneNumber(phone string, otp string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Printf("AWS config error: %v", err)
		return "", err
	}

	client := sns.NewFromConfig(cfg)

	message := fmt.Sprintf("Your OTP is %s. Valid for 5 minutes.", otp)

	output, err := client.Publish(context.Background(), &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(phone),
	})

	if err != nil {
		log.Printf("SNS Publish failed: %v", err)
		return "", err
	}

	log.Printf("SNS MessageID: %s", *output.MessageId)
	return "otp sent", nil
}
