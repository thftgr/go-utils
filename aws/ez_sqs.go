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
	Client   *sqs.Client
	QueueUrl *string
	Timeout  time.Duration
}

// NewEzSqs
//
//	param:
//		cfg *aws.Config          : aws config
//		queueName string         : sqs queue name
//		maxTimeout time.Duration : sqs timeout
func NewEzSqs(cfg *aws.Config, queueName string, maxTimeout time.Duration) (q *EzSqs, err error) {
	q = &EzSqs{Client: sqs.NewFromConfig(*cfg), Timeout: maxTimeout}
	ctx, cancel := context.WithTimeout(context.Background(), q.Timeout)
	defer cancel()
	res, err := q.Client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: &queueName})
	if err != nil {
		return nil, err
	}
	q.QueueUrl = res.QueueUrl

	return
}

func (q *EzSqs) SendMessage(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), q.Timeout)
	defer cancel()
	_, err := q.Client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: &message,
		QueueUrl:    q.QueueUrl,
	})
	return err
}

func (q *EzSqs) GetMessage() (*types.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), q.Timeout) // 큐의 롱폴링 max 대기시간이 20초임.
	defer cancel()

	res, err := q.Client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            q.QueueUrl,
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
	ctx, cancel := context.WithTimeout(context.Background(), q.Timeout) // 큐의 롱폴링 max 대기시간이 20초임.
	defer cancel()

	res, err := q.Client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            q.QueueUrl,
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

func (q *EzSqs) DeleteMessage(receiptHandle *string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), q.Timeout)
	defer cancel()
	_, err = q.Client.DeleteMessage(ctx, &sqs.DeleteMessageInput{QueueUrl: q.QueueUrl, ReceiptHandle: receiptHandle})
	return
}
