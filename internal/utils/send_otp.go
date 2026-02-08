package utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SendExistingOTP sends a provided OTP to the given phone number
func SendExistingOTP(phone string, otp string) (string, error) {
	// 1. Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %v", err)
	}

	// 2. Create SNS client
	snsClient := sns.NewFromConfig(cfg)

	// 3. Prepare message
	message := fmt.Sprintf("Your KHEL app OTP is: %s. It is valid for 5 minutes.", otp)

	// 4. Publish SMS
	_, err = snsClient.Publish(context.TODO(), &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(phone), // Must be in E.164 format e.g., +923001234567
	})
	if err != nil {
		return "", fmt.Errorf("failed to send SMS: %v", err)
	}

	return "otp sent", nil
}
