package queue

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func StartWorkerPool(ctx context.Context, client *sqs.Client, queueURL string, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			for {
				select {
				case <-ctx.Done():
					log.Printf("Worker %d shutting down", workerID)
					return
				default:
					job, receipt, err := Receive(ctx, client, queueURL)
					if err != nil {
						if err == ErrNoMessages {
							continue
						}
						log.Println("Receive error:", err)
						continue
					}

					// Process job
					log.Println("Worker", workerID, "processing job:", job.JobType)
					// call some service to handle job
					_ = Delete(ctx, client, queueURL, receipt)
				}
			}
		}(i)
	}
}
