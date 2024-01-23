package aws

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

const _S3_TEST_AWS_PROFILE = "dev"
const _S3_TEST_BUCKET_NAME = "go-s3-wrapper-test"

func TestEzS3_GetObject(t *testing.T) {
	cfg := NewAwsConfigWUsingProfile(_S3_TEST_AWS_PROFILE)
	tests := []struct {
		name  string
		ezaws *EzS3
		key   string
		body  []byte
	}{
		{name: "", ezaws: NewEzS3(cfg, _S3_TEST_BUCKET_NAME, time.Second*10), key: "test/test.txt", body: []byte("hello?")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tt.ezaws.PutObject(tt.key, bytes.NewReader(tt.body)); err != nil {
				t.Errorf("failed to setup object. err:%+v", err)
				return
			}

			defer func() {
				if err := tt.ezaws.DeleteObject(tt.key); err != nil {
					t.Errorf("failed to delete setup Object(%s)", tt.key)
				}
			}()
			resReader, err := tt.ezaws.GetObject(tt.key)
			if err != nil {
				t.Errorf("GetObject() error = %v", err)
				return
			}
			if body, err := io.ReadAll(resReader); err != nil {
				t.Errorf("GetObject() read body error = %v", err)
				return
			} else if !reflect.DeepEqual(body, tt.body) {
				t.Errorf("GetObject() gotRes = %v, want %v", body, tt.body)
			}
		})
	}
}

func TestEzS3_GetPreSigned(t *testing.T) {
	cfg := NewAwsConfigWUsingProfile(_S3_TEST_AWS_PROFILE)
	tests := []struct {
		name  string
		ezaws *EzS3
		key   string
		body  []byte
	}{
		{name: "", ezaws: NewEzS3(cfg, _S3_TEST_BUCKET_NAME, time.Second*10), key: "test/test.txt", body: []byte("hello?")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tt.ezaws.PutObject(tt.key, bytes.NewReader(tt.body)); err != nil {
				t.Errorf("failed to setup object. err:%+v", err)
				return
			}

			defer func() {
				if err := tt.ezaws.DeleteObject(tt.key); err != nil {
					t.Errorf("failed to delete setup Object(%s)", tt.key)
				}
			}()
			url, err := tt.ezaws.GetPreSigned(tt.key, time.Second*10)
			if err != nil {
				t.Errorf("GetPreSigned() error = %v", err)
				return
			}
			res, err := http.Get(url)
			if err != nil {
				t.Errorf("GetPreSigned() error = %v", err)
				return
			}
			defer res.Body.Close()
			if body, err := io.ReadAll(res.Body); err != nil {
				t.Errorf("GetPreSigned() read body error = %v", err)
				return
			} else if !reflect.DeepEqual(body, tt.body) {
				t.Errorf("GetPreSigned() gotRes = %v, want %v", body, tt.body)
			}
		})
	}
}

func TestEzS3_PutObject(t *testing.T) {
	cfg := NewAwsConfigWUsingProfile(_S3_TEST_AWS_PROFILE)
	tests := []struct {
		name  string
		ezaws *EzS3
		key   string
		body  []byte
	}{
		{name: "", ezaws: NewEzS3(cfg, _S3_TEST_BUCKET_NAME, time.Second*10), key: "test/test.txt", body: []byte("hello?")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tt.ezaws.PutObject(tt.key, bytes.NewReader(tt.body)); err != nil {
				t.Errorf("failed to setup object. err:%+v", err)
				return
			}

			defer func() {
				if err := tt.ezaws.DeleteObject(tt.key); err != nil {
					t.Errorf("failed to delete setup Object(%s)", tt.key)
				}
			}()
			resReader, err := tt.ezaws.GetObject(tt.key)
			if err != nil {
				t.Errorf("PutObject() error = %v", err)
				return
			}
			if body, err := io.ReadAll(resReader); err != nil {
				t.Errorf("PutObject() read body error = %v", err)
				return
			} else if !reflect.DeepEqual(body, tt.body) {
				t.Errorf("PutObject() gotRes = %v, want %v", body, tt.body)
			}
		})
	}
}

func TestEzS3_PutPreSigned(t *testing.T) {
	cfg := NewAwsConfigWUsingProfile(_S3_TEST_AWS_PROFILE)
	tests := []struct {
		name  string
		ezaws *EzS3
		key   string
		body  []byte
	}{
		{name: "", ezaws: NewEzS3(cfg, _S3_TEST_BUCKET_NAME, time.Second*10), key: "test/test.txt", body: []byte("hello?")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := tt.ezaws.PutPreSigned(tt.key, time.Second*10)
			if err != nil {
				t.Errorf("PutPreSigned() error = %v", err)
				return
			}

			defer func() {
				if err := tt.ezaws.DeleteObject(tt.key); err != nil {
					t.Errorf("failed to delete setup Object(%s)", tt.key)
				}
			}()
			req, err := http.NewRequest("PUT", url, bytes.NewReader(tt.body))
			if err != nil {
				t.Errorf("PutPreSigned() error = %v", err)
				return
			}
			client := http.Client{Timeout: time.Second * 10}

			res, err := client.Do(req)
			if err != nil {
				t.Errorf("PutPreSigned() error = %v", err)
				return
			}

			defer res.Body.Close()

			resReader, err := tt.ezaws.GetObject(tt.key)
			if err != nil {
				t.Errorf("PutPreSigned() error = %v", err)
				return
			}
			if body, err := io.ReadAll(resReader); err != nil {
				t.Errorf("PutPreSigned() read body error = %v", err)
				return
			} else if !reflect.DeepEqual(body, tt.body) {
				t.Errorf("PutPreSigned() gotRes = %v, want %v", body, tt.body)
			}
		})
	}
}
