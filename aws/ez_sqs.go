package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"time"
)

// EzSqs 간단한 sqs 작업을 지원합니다.
type EzSqs struct {
	client   *sqs.Client
	queueUrl *string
	timeout  time.Duration
}

func NewEzSqs(cfg *aws.Config, queueName string, maxTimeout time.Duration) (q *EzSqs, err error) {
	q = &EzSqs{client: sqs.NewFromConfig(*cfg), timeout: maxTimeout}
	ctx, cancel := context.WithTimeout(context.Background(), q.timeout)
	defer cancel()
	res, err := q.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: &queueName})
	if err != nil {
		return nil, err
	}
	q.queueUrl = res.QueueUrl

	return
}

func (q *EzSqs) SendMessage(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), q.timeout)
	defer cancel()
	_, err := q.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: &message,
		QueueUrl:    q.queueUrl,
	})
	return err
}

func (q *EzSqs) GetMessage() (*types.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), q.timeout) // 큐의 롱폴링 max 대기시간이 20초임.
	defer cancel()

	res, err := q.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            q.queueUrl,
		MaxNumberOfMessages: 1, // 1개만 가져옴
	})
	if err != nil {
		return nil, err
	}
	if len(res.Messages) == 0 {
		return nil, nil
	}
	return &res.Messages[0], nil

}

// GetMessages
// size range 1~10
// if 0 apply default value 1
func (q *EzSqs) GetMessages(size int) ([]types.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), q.timeout) // 큐의 롱폴링 max 대기시간이 20초임.
	defer cancel()

	res, err := q.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            q.queueUrl,
		MaxNumberOfMessages: int32(size), // n개 가져옴
	})
	if err != nil {
		return nil, err
	}
	if len(res.Messages) == 0 {
		return nil, nil
	}
	return res.Messages, nil
}
