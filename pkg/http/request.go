package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func NewRequest(c *Client) *Request {
	return &Request{
		c:       c,
		baseURL: c.baseURL,
		params:  make(url.Values),
		headers: make(http.Header),
	}
}

type Request struct {
	c       *Client
	baseURL string
	subURL  string
	params  url.Values
	headers http.Header
	method  string
	body    io.Reader
}

type Result struct {
	Body       []byte
	Err        error
	StatusCode int
}

func (r *Request) SetBaseURL(baseURL string) *Request {
	r.baseURL = baseURL
	return r
}

func (r *Request) SetSubURL(subURL string) *Request {
	r.subURL = subURL
	return r
}

func (r *Request) AddParamsFromValues(params url.Values) *Request {
	for k, vv := range params {
		for _, v := range vv {
			r.params.Add(k, v)
		}
	}
	return r
}

func (r *Request) SetParams(params map[string]string) *Request {
	for p, v := range params {
		r.SetParam(p, v)
	}
	return r
}

func (r *Request) SetParam(param, value string) *Request {
	r.params.Set(param, value)
	return r
}

func (r *Request) AddHeaders(headers http.Header) *Request {
	for k, hh := range headers {
		for _, h := range hh {
			r.headers.Add(k, h)
		}
	}
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	for h, v := range headers {
		r.SetHeader(h, v)
	}
	return r
}

func (r *Request) SetHeader(header, value string) *Request {
	r.headers.Set(header, value)
	return r
}

func (r *Request) SetMethod(method string) *Request {
	r.method = method
	return r
}

func (r *Request) SetBody(body io.Reader) *Request {
	r.body = body
	return r
}

func (r *Request) URL() (*url.URL, error) {
	finalURL, err := url.Parse(r.baseURL)
	if err != nil {
		return nil, err
	}
	finalURL.Path = r.subURL
	query := url.Values{}
	for key, values := range r.params {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	finalURL.RawQuery = query.Encode()
	return finalURL, nil
}

func (r *Request) DoRequest(ctx context.Context) *Result {
	finalURL, err := r.URL()
	if err != nil {
		return &Result{
			Err: fmt.Errorf("baseURL is not correct format: %w", err),
		}
	}
	req, err := http.NewRequestWithContext(ctx, r.method, finalURL.String(), r.body)
	if err != nil {
		return &Result{
			Err: fmt.Errorf("failed to make request: %w", err),
		}
	}
	req.Header = r.headers
	res, err := r.c.httpClient.Do(req)
	if err != nil {
		return &Result{
			Err: fmt.Errorf("failed to do request: %w", err),
		}
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &Result{
			Err:        fmt.Errorf("failed to read response body: %w", err),
			StatusCode: res.StatusCode,
		}
	}
	return &Result{
		Body:       body,
		Err:        nil,
		StatusCode: res.StatusCode,
	}
}
