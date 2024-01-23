package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"reflect"
	"testing"
	"time"
)

func TestEzSqs_GetMessage(t *testing.T) {
	tests := []struct {
		name     string
		client   *sqs.Client
		queueUrl *string
		timeout  time.Duration
		want     *types.Message
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &EzSqs{
				Client:   tt.client,
				QueueUrl: tt.queueUrl,
				Timeout:  tt.timeout,
			}
			got, err := q.GetMessage()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEzSqs_GetMessages(t *testing.T) {
	tests := []struct {
		name     string
		client   *sqs.Client
		queueUrl *string
		timeout  time.Duration
		size     int
		want     []types.Message
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &EzSqs{
				Client:   tt.client,
				QueueUrl: tt.queueUrl,
				Timeout:  tt.timeout,
			}
			got, err := q.GetMessages(tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMessages() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEzSqs_SendMessage(t *testing.T) {
	tests := []struct {
		name     string
		client   *sqs.Client
		queueUrl *string
		timeout  time.Duration
		message  string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &EzSqs{
				Client:   tt.client,
				QueueUrl: tt.queueUrl,
				Timeout:  tt.timeout,
			}
			if err := q.SendMessage(tt.message); (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewEzSqs(t *testing.T) {
	tests := []struct {
		name       string
		cfg        *aws.Config
		queueName  string
		maxTimeout time.Duration
		wantQ      *EzSqs
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQ, err := NewEzSqs(tt.cfg, tt.queueName, tt.maxTimeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEzSqs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotQ, tt.wantQ) {
				t.Errorf("NewEzSqs() gotQ = %v, want %v", gotQ, tt.wantQ)
			}
		})
	}
}
