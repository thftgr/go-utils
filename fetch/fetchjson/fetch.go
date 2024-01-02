package fetchjson

import (
	"errors"
	"io"
	"net/http"
	"time"
)

type JsonFetch struct {
	Header http.Header
	Client *http.Client
}

func NewJsonFetch(header http.Header, timeout time.Duration) *JsonFetch {
	return &JsonFetch{
		Header: header,
		Client: &http.Client{Timeout: timeout},
	}
}

func (r *JsonFetch) parse(resp *http.Response, ei error) (responseBody []byte, eo error) {
	defer func() {
		if resp != nil && !resp.Close {
			eo = errors.Join(eo, resp.Body.Close())
		}
	}()
	if ei != nil {
		eo = ei
		return
	}

	responseBody, eo = io.ReadAll(resp.Body)

	if 200 <= resp.StatusCode && resp.StatusCode < 300 {
		return responseBody, eo
	} else {
		return responseBody, errors.Join(eo, errors.New(resp.Status))
	}
}

func (r *JsonFetch) Get(url string) (responseBody []byte, err error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	for k := range r.Header {
		request.Header.Add(k, r.Header.Get(k))
	}
	return r.parse(r.Client.Do(request))
}
func (r *JsonFetch) Head(url string) (responseBody []byte, err error) {
	request, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return
	}
	for k := range r.Header {
		request.Header.Add(k, r.Header.Get(k))
	}
	return r.parse(r.Client.Do(request))
}
func (r *JsonFetch) Post(url string, requestBody io.ReadCloser) (responseBody []byte, err error) {
	defer func() {
		if requestBody != nil {
			err = errors.Join(err, requestBody.Close())
		}
	}()
	request, err := http.NewRequest(http.MethodPost, url, requestBody)
	if err != nil {
		return
	}
	for k := range r.Header {
		request.Header.Add(k, r.Header.Get(k))
	}
	return r.parse(r.Client.Do(request))
}
func (r *JsonFetch) PUT(url string, requestBody io.ReadCloser) (responseBody []byte, err error) {
	defer func() {
		if requestBody != nil {
			err = errors.Join(err, requestBody.Close())
		}
	}()
	request, err := http.NewRequest(http.MethodPut, url, requestBody)
	if err != nil {
		return
	}
	for k := range r.Header {
		request.Header.Add(k, r.Header.Get(k))
	}
	return r.parse(r.Client.Do(request))
}
func (r *JsonFetch) PATCH(url string, requestBody io.ReadCloser) (responseBody []byte, err error) {
	defer func() {
		if requestBody != nil {
			err = errors.Join(err, requestBody.Close())
		}
	}()
	request, err := http.NewRequest(http.MethodPatch, url, requestBody)
	if err != nil {
		return
	}
	for k := range r.Header {
		request.Header.Add(k, r.Header.Get(k))
	}
	return r.parse(r.Client.Do(request))
}
func (r *JsonFetch) DELETE(url string) (responseBody []byte, err error) {
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return
	}
	for k := range r.Header {
		request.Header.Add(k, r.Header.Get(k))
	}
	return r.parse(r.Client.Do(request))
}
func (r *JsonFetch) CONNECT(url string) (responseBody []byte, err error) {
	request, err := http.NewRequest(http.MethodConnect, url, nil)
	if err != nil {
		return
	}
	for k := range r.Header {
		request.Header.Add(k, r.Header.Get(k))
	}
	return r.parse(r.Client.Do(request))
}
func (r *JsonFetch) OPTIONS(url string) (responseBody []byte, err error) {
	request, err := http.NewRequest(http.MethodOptions, url, nil)
	if err != nil {
		return
	}
	for k := range r.Header {
		request.Header.Add(k, r.Header.Get(k))
	}
	return r.parse(r.Client.Do(request))
}
func (r *JsonFetch) TRACE(url string) (responseBody []byte, err error) {
	request, err := http.NewRequest(http.MethodTrace, url, nil)
	if err != nil {
		return
	}
	for k := range r.Header {
		request.Header.Add(k, r.Header.Get(k))
	}
	return r.parse(r.Client.Do(request))
}

//	MethodGet     = "GET"
//	MethodHead    = "HEAD"
//	MethodPost    = "POST"
//	MethodPut     = "PUT"
//	MethodPatch   = "PATCH" // RFC 5789
//	MethodDelete  = "DELETE"
//	MethodConnect = "CONNECT"
//	MethodOptions = "OPTIONS"
//	MethodTrace   = "TRACE"
