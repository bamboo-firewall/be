package http

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

// mockReader mock for body reader
type mockReader struct{}

// Read interface implement io.Reader and always return error
func (m mockReader) Read(p []byte) (int, error) {
	return 0, errors.New("simulated read error")
}

func (m mockReader) Close() error {
	return nil
}

func TestDoRequest(t *testing.T) {
	c := NewClient("")
	tests := []struct {
		name                 string
		ctx                  context.Context
		baseURL              string
		subURL               string
		addParams            url.Values
		setParams            map[string]string
		addHeaders           http.Header
		setHeaders           map[string]string
		method               string
		body                 string
		readResponseBodyFail bool
		err                  error
		statusCode           int
	}{
		{
			name:    "parse base url fail",
			baseURL: ":http://localhost",
			err:     errors.New("parse base url fail"),
		},
		{
			name: "new request fail",
			err:  errors.New("new request fail"),
		},
		{
			name:    "do request fail",
			ctx:     context.Background(),
			baseURL: "http://localhost.xyz",
			err:     errors.New("do request fail"),
		},
		{
			name:                 "do request fail",
			ctx:                  context.Background(),
			readResponseBodyFail: true,
			err:                  errors.New("read response body fail"),
		},
		{
			name: "new request success with setParams",
			ctx:  context.Background(),
			setParams: map[string]string{
				"abc": "xyz",
			},
			err:        nil,
			statusCode: http.StatusOK,
		},
		{
			name: "new request success with setHeaders",
			ctx:  context.Background(),
			setHeaders: map[string]string{
				"AccessToken": "aaaabbbbb",
			},
			err:        nil,
			statusCode: http.StatusOK,
		},
		{
			name: "new request success with all",
			ctx:  context.Background(),
			addParams: url.Values{
				"ids": []string{"1", "2", "3"},
				"abc": []string{"xyz"},
			},
			setParams: map[string]string{
				"abc": "xyz",
			},
			addHeaders: http.Header{
				"header1": []string{"value1"},
			},
			setHeaders: map[string]string{
				"AccessToken": "aaaabbbbb",
				"header1":     "value2",
			},
			method:     http.MethodPost,
			body:       `{"aaaa":"bbbb"}`,
			statusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestParams := make(url.Values)
			if tt.addParams != nil {
				for k, vv := range tt.addParams {
					for _, v := range vv {
						requestParams.Add(k, v)
					}
				}
			}
			for k, v := range tt.setParams {
				requestParams.Set(k, v)
			}

			requestHeaders := make(http.Header)
			if tt.addHeaders != nil {
				for k, vv := range tt.addHeaders {
					for _, v := range vv {
						requestHeaders.Add(k, v)
					}
				}
			}
			for k, v := range tt.setHeaders {
				requestHeaders.Set(k, v)
			}

			var bodyReader io.Reader
			if tt.body != "" {
				bodyReader = bytes.NewBufferString(tt.body)
			}
			handler := func(w http.ResponseWriter, r *http.Request) {
				if tt.readResponseBodyFail {
					w.Header().Set("Content-Length", "1")
					return
				}
				if len(requestParams) > 0 {
					for k, v := range requestParams {
						if !cmp.Equal(v, r.URL.Query()[k]) {
							w.WriteHeader(http.StatusBadRequest)
							return
						}
					}
				}

				if len(requestHeaders) > 0 {
					for k, h := range requestHeaders {
						if !cmp.Equal(h, r.Header[k]) {
							w.WriteHeader(http.StatusBadRequest)
							return
						}
					}
				}

				if tt.body != "" {
					body, err := io.ReadAll(r.Body)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					if string(body) != tt.body {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
				}
			}

			server := httptest.NewServer(http.HandlerFunc(handler))
			defer server.Close()
			if tt.baseURL == "" {
				tt.baseURL = server.URL
			}

			result := c.NewRequest().
				SetBaseURL(tt.baseURL).
				SetSubURL(tt.subURL).
				AddParamsFromValues(tt.addParams).
				SetParams(tt.setParams).
				AddHeaders(tt.addHeaders).
				SetHeaders(tt.setHeaders).
				SetMethod(tt.method).
				SetBody(bodyReader).
				DoRequest(tt.ctx)
			if tt.err != nil {
				assert.Error(t, result.Err)
				return
			}
			assert.Equal(t, tt.statusCode, result.StatusCode)
		})
	}
}
