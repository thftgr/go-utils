package aws

import (
	"reflect"
	"testing"
	"time"
)

const _SQS_TEST_AWS_PROFILE = "dev"
const _SQS_TEST_QUEUE_NAME = "go-sqs-wrapper-test"

// TestEzSqs_SendMessage_GetMessage_DeleteMessage
// 중복되는 테스트가 많이 통합
func TestEzSqs_SendMessage_GetMessage_DeleteMessage(t *testing.T) {
	ezsqs, err := NewEzSqs(NewAwsConfigWUsingProfile(_SQS_TEST_AWS_PROFILE), _SQS_TEST_QUEUE_NAME, time.Second*30)
	if err != nil {
		t.Errorf("GetMessage() error = %v", err)
		return
	}
	tests := []struct {
		name    string
		ezsqs   *EzSqs
		message string
	}{
		{name: "", ezsqs: ezsqs, message: "hello?"},
		{name: "", ezsqs: ezsqs, message: "world!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ezsqs.SendMessage(tt.message); err != nil {
				t.Errorf("SendMessage() error = %v", err)
				return
			}

			msg, err := tt.ezsqs.GetMessage()
			if err != nil {
				t.Errorf("GetMessage() error = %v", err)
				return
			}
			t.Logf("GetMessage() message:%v", *msg.Body)
			defer func() {
				if err := tt.ezsqs.DeleteMessage(msg.ReceiptHandle); err != nil {
					t.Errorf("DeleteMessage() error = %v", err)
					return
				}
				if msg, err := tt.ezsqs.GetMessage(); err != nil {
					t.Errorf("GetMessage() error = %v", err)
					return
				} else if msg != nil {
					t.Errorf("DeleteMessage() is not deleted : msg:%v", msg.Body)
					return
				}

			}()
			if !reflect.DeepEqual(*msg.Body, tt.message) {
				t.Errorf("GetMessage() got = %v, message %v", *msg.Body, tt.message)
			}
		})
	}
}

// TestEzSqs_SendMessage_GetMessage_DeleteMessage
// 중복되는 테스트가 많이 통합
func TestEzSqs_SendMessage_GetMessages_DeleteMessage(t *testing.T) {
	ezsqs, err := NewEzSqs(NewAwsConfigWUsingProfile(_SQS_TEST_AWS_PROFILE), _SQS_TEST_QUEUE_NAME, time.Second*30)
	if err != nil {
		t.Errorf("GetMessage() error = %v", err)
		return
	}
	tests := []struct {
		name    string
		ezsqs   *EzSqs
		message []string
	}{
		{name: "", ezsqs: ezsqs, message: []string{"hello!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, s := range tt.message {
				if err := tt.ezsqs.SendMessage(s); err != nil {
					t.Errorf("SendMessage() error = %v", err)
					return
				}
			}

			msg, err := tt.ezsqs.GetMessages(len(tt.message))
			if err != nil {
				t.Errorf("GetMessage() error = %v", err)
				return
			}
			defer func() {
				for _, message := range msg {
					if err := tt.ezsqs.DeleteMessage(message.ReceiptHandle); err != nil {
						t.Errorf("DeleteMessage() error = %v", err)
					}
				}

				if msg, err := tt.ezsqs.GetMessages(len(tt.message)); err != nil {
					t.Errorf("GetMessage() error = %v", err)
					return

				} else if msg != nil {
					for _, message := range msg {
						t.Errorf("DeleteMessage() is not deleted : msg: %v", *message.Body)
					}
					return
				}
			}()
			for i, message := range msg {
				if !reflect.DeepEqual(*message.Body, tt.message[i]) {
					t.Errorf("GetMessage() got = %v, message %v", *message.Body, tt.message[i])
				}
			}

		})
	}
}
