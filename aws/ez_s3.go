package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"time"
)

// EzS3 간단한 s3 작업을 지원합니다.
type EzS3 struct {
	Client        *s3.Client
	Bucket        *string
	Timeout       time.Duration
	PresignClient *s3.PresignClient
}

// NewEzS3
//
//	param:
//		cfg *aws.Config      : aws config
//		bucket string        : bucket name
//		timeout time.Duration: timeout ex) time.Second * 60
func NewEzS3(cfg *aws.Config, bucket string, timeout time.Duration) (s *EzS3) {
	s = &EzS3{Client: s3.NewFromConfig(*cfg), Bucket: &bucket, Timeout: timeout}
	s.PresignClient = s3.NewPresignClient(s.Client)
	return
}

// GetObject
//
// param:
//
//	O key string: my/bucket/path/filename.png
//	X key string: /my/bucket/path/filename.png
//	X key string: s3://mybucket/my/bucket/path/filename.png
//
// return:
//
//	res io.ReadCloser : io.EOF 가 발생할때까지 읽고, 읽음 여부와 상관 없이 반드시 닫아야합니다.
//	err error         : error
func (b *EzS3) GetObject(key string) (res io.ReadCloser, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	result, err := b.Client.GetObject(ctx, &s3.GetObjectInput{Bucket: b.Bucket, Key: &key})
	if err != nil {
		return nil, err
	}
	// object 를 미리 읽어서 []byte 로 반환하는경우 크기에 따라 OOM 발생할수있음.
	return result.Body, err

}

// GetPreSigned
// param:
//
//	O key string: my/bucket/path/filename.png
//	X key string: /my/bucket/path/filename.png
//	X key string: s3://mybucket/my/bucket/path/filename.png
//
//	expires time.Duration: PreSigned url의 만료시간
//
// return:
//
//	url string : PreSigned url
//	err error  : error
func (b *EzS3) GetPreSigned(key string, expires time.Duration) (url string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	res, err := b.PresignClient.PresignGetObject(
		ctx,
		&s3.GetObjectInput{Bucket: b.Bucket, Key: &key},
		func(options *s3.PresignOptions) {
			options.Expires = expires
		})
	if err != nil {
		return "", err
	} else {
		return res.URL, nil
	}
}

// PutPreSigned
// param:
//
//	O key string: my/bucket/path/filename.png
//	X key string: /my/bucket/path/filename.png
//	X key string: s3://mybucket/my/bucket/path/filename.png
//
//	expires time.Duration: PreSigned url의 만료시간
//
// return:
//
//	url string : PreSigned url
//	err error  : error
func (b *EzS3) PutPreSigned(key string, expires time.Duration) (url string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	res, err := b.PresignClient.PresignPutObject(
		ctx,
		&s3.PutObjectInput{Bucket: b.Bucket, Key: &key},
		func(options *s3.PresignOptions) {
			options.Expires = expires
		})
	if err != nil {
		return "", err
	} else {
		return res.URL, nil
	}
}

// PutObject
// param:
//
//	O key string: my/bucket/path/filename.png
//	X key string: /my/bucket/path/filename.png
//	X key string: s3://mybucket/my/bucket/path/filename.png
//
//	body io.Reader: object body
//
// return:
//
//	error         : error
func (b *EzS3) PutObject(key string, body io.Reader) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	_, err = b.Client.PutObject(ctx, &s3.PutObjectInput{Bucket: b.Bucket, Key: &key, Body: body})
	return
}

// DeleteObject
// param:
//
//	O key string: my/bucket/path/filename.png
//	X key string: /my/bucket/path/filename.png
//	X key string: s3://mybucket/my/bucket/path/filename.png
//
// return:
//
//	error         : error
func (b *EzS3) DeleteObject(key string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	_, err = b.Client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: b.Bucket, Key: &key})
	return
}
