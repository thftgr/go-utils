package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"sync"
	"time"
)

// SqsSender 기본적인 기능만 구현되어있음 필요한경우 embed 하여 구현
type SqsSender struct {
	Client   *sqs.Client
	QueueUrl string
}

func (s *SqsSender) SendMessage(ctx context.Context, message string) error {
	_, err := s.Client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: &message,
		QueueUrl:    &s.QueueUrl,
	})
	return err
}

// AsyncSqsReceiver
// 10개의 메세지를 한번에 가져오고 비동기로 처리함
type AsyncSqsReceiver struct {
	Client          *sqs.Client
	QueueUrl        string
	Timeout         time.Duration
	WaitTimeSeconds int32
	ErrorHandler    func(error)
	HandleFunc      func(types.Message) error

	//for worker
	wg   sync.WaitGroup
	pool chan types.Message
}

func (s *AsyncSqsReceiver) Start(ctx context.Context) error {
	s.pool = make(chan types.Message, 10) // sqs 크기가 10임.
	s.runWorker(ctx, 10)

	go func() {
		rmi := &sqs.ReceiveMessageInput{
			QueueUrl:            &s.QueueUrl,
			MaxNumberOfMessages: 10,
			//VisibilityTimeout:       0,
			WaitTimeSeconds: 20,
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
				rctx, cancel := context.WithTimeout(ctx, s.Timeout)
				res, err := s.Client.ReceiveMessage(rctx, rmi) // 메세지를 10개씩 가져와서 pool 에 추가함
				cancel()
				if err != nil && s.ErrorHandler != nil {
					s.ErrorHandler(err)
				}
				for i := range res.Messages {
					s.pool <- res.Messages[i] // pool size 가 넘어가면 blocking 됨.
				}
			}

		}
	}()
	return nil
}

func (s *AsyncSqsReceiver) runWorker(c context.Context, size int) {
	for i := 0; i < size; i++ {
		go func() {
			s.wg.Add(1)
			defer s.wg.Done()
			for {
				select {
				case <-c.Done():
					return
				case msg := <-s.pool:
					if err := s.HandleFunc(msg); err != nil {
						s.ErrorHandler(err)
					}
				}
			}
		}()
	}
	return
}
