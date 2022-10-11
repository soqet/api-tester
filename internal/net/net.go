package net

import (
	"api-tester/internal/jsonreader"
	"bytes"
	"net"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"
)

func Exec(requests *jsonreader.MainSchema) ([]Info, error) {
	results, err := execRequests(requests.Requests, requests.Endpoint, requests.Headers)
	if err != nil {
		return nil, err
	}
	for _, g := range requests.Groups {
		endpoint := requests.Endpoint
		if g.Endpoint != "" {
			endpoint = g.Endpoint
		}
		addRes, err := execRequests(g.Requests, endpoint, g.Headers)
		if err != nil {
			return nil, err
		}
		results = append(results, addRes...)
	}
	return results, nil
}

func execRequests(requests []*jsonreader.RequestSchema, baseEndpoint string, baseHeaders []*jsonreader.HeaderSchema) ([]Info, error) {
	result := make([]Info, len(requests))
	for i := 0; i < len(requests); i++ {
		result[i].Passed = true
	}
	for i, req := range requests {
		body := bytes.NewReader([]byte(req.Body))
		endpoint := baseEndpoint
		if req.Endpoint != "" {
			endpoint = req.Endpoint
		}
		r, err := http.NewRequest(req.Method, endpoint+req.Resourse, body)
		if err != nil {
			return nil, err
		}
		c := http.Client{
			Timeout: time.Second * 2,
		}
		var start time.Time
		var elapsed int64
		clientTrace := &httptrace.ClientTrace{
			GotConn: func(_ httptrace.GotConnInfo) { start = time.Now() },
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
		result[i].Time = elapsed
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				result[i].Passed = false
				result[i].Reason = "Request timeout"
				continue
			}
			return nil, err
		}
		if resp.StatusCode != req.Code {
			result[i].Passed = false
			result[i].Reason = "Wrong code"
			result[i].Code = resp.StatusCode
			continue
		}
		httpBuf := bytes.Buffer{}
		httpBuf.ReadFrom(resp.Body)
		testResp := []byte(req.Response)
		if req.ResponseFile != "" {
			testResp, err = readFile(req.ResponseFile)
			if err != nil {
				return nil, err
			}
		}
		if !bytes.Equal(httpBuf.Bytes(), testResp) {
			result[i].Passed = false
			result[i].Reason = "Wrong body"
			result[i].Response = httpBuf.String()
		}
	}
	return result, nil
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