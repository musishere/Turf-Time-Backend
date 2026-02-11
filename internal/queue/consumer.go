package queue

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/musishere/sportsApp/types"
)

var ErrNoMessages = errors.New("no messages in queue")

func Receive(ctx context.Context, client *sqs.Client, queueURL string) (types.Job, string, error) {
	resp, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     20, // long polling
	})
	if err != nil {
		return types.Job{}, "", err
	}

	if len(resp.Messages) == 0 {
		return types.Job{}, "", ErrNoMessages
	}

	msg := resp.Messages[0]
	var job types.Job
	if err := json.Unmarshal([]byte(*msg.Body), &job); err != nil {
		return types.Job{}, "", err
	}

	return job, *msg.ReceiptHandle, nil
}
