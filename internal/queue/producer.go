package queue

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/musishere/sportsApp/types"
)

func Send(ctx context.Context, client *sqs.Client, queueURL string, job types.Job) (string, error) {
	body, err := json.Marshal(job)
	if err != nil {
		return "Error sending the message to the queue", err
	}

	_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(body)),
	})

	return "Message sent to the queue", err
}
