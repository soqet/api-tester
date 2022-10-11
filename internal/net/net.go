package net

import (
	"api-tester/internal/jsonreader"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"
)

func Exec(requests *jsonreader.MainSchema, results chan Info) error {
	defer close(results)
	err := execRequests(requests.Requests, requests.Endpoint, requests.Headers, results)
	if err != nil {
		return err
	}
	for _, g := range requests.Groups {
		endpoint := requests.Endpoint
		if g.Endpoint != "" {
			endpoint = g.Endpoint
		}
		err = execRequests(g.Requests, endpoint, g.Headers, results)
		if err != nil {
			return err
		}
	}
	return nil
}

func execRequests(
	requests []*jsonreader.RequestSchema, baseEndpoint string,
	baseHeaders []*jsonreader.HeaderSchema, results chan Info,
) error {
	timeout := time.Second * 3
	for _, req := range requests {
		body := bytes.NewReader([]byte(req.Body))
		endpoint := baseEndpoint
		if req.Endpoint != "" {
			endpoint = req.Endpoint
		}
		r, err := http.NewRequest(req.Method, endpoint+req.Resourse, body)
		if err != nil {
			return err
		}
		c := http.Client{
			Timeout: timeout,
		}
		var start time.Time
		var elapsed int64
		clientTrace := &httptrace.ClientTrace{
			GotConn:              func(_ httptrace.GotConnInfo) { start = time.Now() },
			GotFirstResponseByte: func() { elapsed = time.Since(start).Milliseconds() },
		}
		r = r.WithContext(httptrace.WithClientTrace(r.Context(), clientTrace))
		for _, h := range append(baseHeaders, req.Headers...) {
			if http.CanonicalHeaderKey(h.Key) == "Host" {
				r.Host = h.Value
			}
			r.Header.Add(h.Key, h.Value)
		}
		resp, err := c.Do(r)
		res := Info{
			Time:   elapsed,
			Passed: true,
		}
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				res.Passed = false
				res.Reason = fmt.Sprintf("Request timeout (%d ms)", timeout/time.Millisecond)
				results <- res
				continue
			}
			return err
		}
		if resp.StatusCode != req.Code {
			res.Passed = false
			res.Reason = "Wrong code"
			res.Code = resp.StatusCode
			results <- res
			continue
		}
		httpBuf := bytes.Buffer{}
		httpBuf.ReadFrom(resp.Body)
		testResp := []byte(req.Response)
		if req.ResponseFile != "" {
			testResp, err = readFile(req.ResponseFile)
			if err != nil {
				return err
			}
		}
		if !bytes.Equal(httpBuf.Bytes(), testResp) {
			res.Passed = false
			res.Reason = "Wrong body"
			res.Response = httpBuf.String()
		}
		results <- res
	}
	return nil
}

func readFile(path string) ([]byte, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	testBuf := bytes.Buffer{}
	testBuf.ReadFrom(file)
	return testBuf.Bytes(), nil
}
