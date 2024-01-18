package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"sync"
	"time"
)

type SqsSender interface {
	SendMessage(message string) error
}

// SqsSenderImpl 기본적인 기능만 구현되어있음 필요한경우 embed 하여 구현
type SqsSenderImpl struct {
	Client     *sqs.Client
	QueueUrl   string
	MaxTimeout time.Duration
}

func NewSqsSenderImpl(client *sqs.Client, queueUrl string, maxTimeout time.Duration) *SqsSenderImpl {
	return &SqsSenderImpl{Client: client, QueueUrl: queueUrl, MaxTimeout: maxTimeout}
}

func (s *SqsSenderImpl) SendMessage(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.MaxTimeout)
	defer cancel()
	_, err := s.Client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: &message,
		QueueUrl:    &s.QueueUrl,
	})

	return err
}

type SqsReceiver interface {
	Start(context.Context) error
}

type AsyncSqsReceiverImpl struct {
	Client          *sqs.Client
	QueueUrl        string
	Timeout         time.Duration
	WaitTimeSeconds int32
	ErrorHandler    func(error)
	HandleFunc      func(types.Message) error

	//for worker
	wg      sync.WaitGroup
	running bool
	pool    chan types.Message
}

func (s *AsyncSqsReceiverImpl) Start(ctx context.Context) error {
	s.pool = make(chan types.Message, 20) // sqs 크기가 10임.

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
				res, err := s.Client.ReceiveMessage(rctx, rmi)
				cancel()
				if err != nil {
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

func (s *AsyncSqsReceiverImpl) runWorker(ctx context.Context, size int) (stop func()) {
	wctx, cancel := context.WithCancel(ctx)
	stop = cancel

	for i := 0; i < size; i++ {
		go func() {
			s.wg.Add(1)
			defer s.wg.Done()

			for {
				select {
				case <-wctx.Done():
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
