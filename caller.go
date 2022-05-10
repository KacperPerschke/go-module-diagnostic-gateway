package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

type myRespType struct {
	Payload     []byte
	ContentType string
}

func callProxy(path string) (myRespType, error) {
	var myResp myRespType
	resp, err := http.Get(fmt.Sprintf("https://proxy.golang.org%s", path))
	defer resp.Body.Close()
	if err != nil {
		return myResp, err
	}
	b, err := io.ReadAll(resp.Body)
	myResp.Payload = b
	myResp.ContentType = resp.Header.Get("Content-Type")
	return myResp, err
}

func (r myRespType) toLog() string {
	if r.ContentType == "application/zip" {
		return base64.StdEncoding.EncodeToString(r.Payload)
	}
	return string(r.Payload)
}
