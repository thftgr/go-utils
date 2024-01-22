package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"reflect"
	"testing"
	"time"
)

const (
	_TEST_AWS_PROFILE    = "dev"
	_TEST_AWS_QUEUE_NAME = "go-sqs-wrapper-test"
)

// Test_SendMessage 테스트 전에 큐를 비우세요.
func Test_SendMessage(t *testing.T) {
	client := sqs.NewFromConfig(*NewAwsConfigWUsingProfile(_TEST_AWS_PROFILE))
	res, err := client.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{QueueName: aws.String(_TEST_AWS_QUEUE_NAME)})
	if err != nil {
		t.Errorf("failed to get queueUrl. error:%+v", err)
		return
	}
	sender := &SqsSender{Client: client, QueueUrl: *res.QueueUrl}
	type testCase struct {
		name     string
		queueUrl *string
		sender   *SqsSender
		timeout  time.Duration
		send     string
		want     string
	}
	tests := []testCase{
		{name: "", queueUrl: res.QueueUrl, sender: sender, timeout: time.Second * 10, send: `hello`, want: `hello`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()
			// 큐 정리

			if err := tt.sender.SendMessage(ctx, tt.send); err != nil {
				t.Errorf("SendMessage() err: %+v", err)
				return
			}

			recv, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{QueueUrl: tt.queueUrl})
			if err != nil {
				t.Errorf("ReceiveMessage() err: %+v", err)
				return
			} else {
				for _, msg := range recv.Messages {
					defer func() {
						client.DeleteMessage(ctx, &sqs.DeleteMessageInput{QueueUrl: res.QueueUrl, ReceiptHandle: msg.ReceiptHandle})
					}()
				}
			}

			if len(recv.Messages) != 1 {
				t.Errorf("ReceiveMessage() receive message != 1")
				return
			} else if len(recv.Messages) < 1 {
				t.Errorf("ReceiveMessage() receive message is empty")
				return
			} else if !reflect.DeepEqual(recv.Messages[0].Body, &tt.want) {
				t.Errorf("Add() want = %s, get %s", tt.want, *recv.Messages[0].Body)
				return
			}

		})
	}

}
