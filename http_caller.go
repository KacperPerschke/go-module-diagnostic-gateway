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
	resp, err := http.Get(fmt.Sprintf("https://%s%s", host, path))
	defer resp.Body.Close()
	if err != nil {
		return myRespType{}, err
	}
	b, err := io.ReadAll(resp.Body)
	myResp := myRespType{
		Header:      resp.Header,
		Payload:     b,
		Status:      resp.Status,
		StatusCode:  resp.StatusCode,
		LogAsBase64: false,
	}
	return myResp, err
}

func (r myRespType) toLog() string {
	if r.LogAsBase64 {
		return base64.StdEncoding.EncodeToString(r.Payload)
	}
	return string(r.Payload)
}
