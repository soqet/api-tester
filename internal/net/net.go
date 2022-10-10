package net

import (
	"api-tester/internal/jsonreader"
	"bytes"
	"net/http"
	"os"
	"time"
)

func Exec(requests *jsonreader.MainSchema) ([]Info, error) {
	results, err := execRequests(requests.Requests, requests.Endpoint)
	if err != nil {
		return nil, err
	}
	for _, g := range requests.Groups {
		endpoint := requests.Endpoint
		if g.Endpoint != "" {
			endpoint = g.Endpoint
		}
		addRes, err := execRequests(g.Requests, endpoint)
		if err != nil {
			return nil, err
		}
		results = append(results, addRes...)
	}
	return results, nil
}

func execRequests(requests []jsonreader.RequestSchema, baseEndpoint string) ([]Info, error) {
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
		c := http.Client{}
		start := time.Now()
		resp, err := c.Do(r)
		elapsed := time.Since(start).Milliseconds()
		result[i].Time = elapsed
		if err != nil {
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
			file, err := os.OpenFile(req.ResponseFile, os.O_RDONLY, 0666)
			if err != nil {
				return nil, err
			}
			testBuf := bytes.Buffer{}
			testBuf.ReadFrom(file)
			file.Close()
			testResp = testBuf.Bytes()
		}
		if !bytes.Equal(httpBuf.Bytes(), testResp) {
			result[i].Passed = false
			result[i].Reason = "Wrong body"
			result[i].Response = httpBuf.String()
		}
	}
	return result, nil
}
