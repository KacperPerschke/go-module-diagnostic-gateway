package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

type myRespType struct {
	Header            http.Header
	Payload           []byte
	Status            string
	StatusCode        int
	ContentTypeToBase bool
	SumDBTile         bool
	LogAsBase64       bool
}

func callWorld(host string, path string) (myRespType, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", host, path))
	if err != nil {
		return myRespType{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	myResp := myRespType{
		Header:      resp.Header.Clone(),
		Payload:     b,
		Status:      resp.Status,
		StatusCode:  resp.StatusCode,
		LogAsBase64: false,
	}
	return myResp, err
}

func (r myRespType) toLog(logAsBase64 bool) string {
	if logAsBase64 {
		return base64.StdEncoding.EncodeToString(r.Payload)
	}
	return string(r.Payload)
}

func (r myRespType) copyHeadersTo(w http.ResponseWriter) {
	for n, values := range r.Header {
		for _, v := range values {
			w.Header().Add(n, v)
		}
	}
}
