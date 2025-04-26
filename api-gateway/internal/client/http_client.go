package client

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// Forward proxies the incoming Gin context to another service and writes the response back.
func Forward(method, url string, body []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return httpClient.Do(req)
}

// CopyResponse copies status, headers and body from src to dst.
func CopyResponse(dst http.ResponseWriter, src *http.Response) error {
	for k, vv := range src.Header {
		for _, v := range vv {
			dst.Header().Add(k, v)
		}
	}
	dst.WriteHeader(src.StatusCode)
	_, err := io.Copy(dst, src.Body)
	return err
}
