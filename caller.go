package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type myRespType struct {
	Header            http.Header
	Payload           []byte
	Status            string
	StatusCode        int
	ContentTypeToBase bool
	SumDBTile         bool
}

func callWorld(host string, path string) (myRespType, error) {
	var myResp myRespType
	resp, err := http.Get(fmt.Sprintf("https://%s%s", host, path))
	defer resp.Body.Close()
	if err != nil {
		return myResp, err
	}
	myResp.Header = resp.Header
	b, err := io.ReadAll(resp.Body)
	myResp.Payload = b
	myResp.Status = resp.Status
	myResp.StatusCode = resp.StatusCode
	myResp.ContentTypeToBase = resp.Header.Get("Content-Type") == "application/zip"
	myResp.SumDBTile = strings.Contains(
		resp.Request.URL.Path,
		`tile`,
	)
	return myResp, err
}

func (r myRespType) toLog() string {
	if r.ContentTypeToBase || r.SumDBTile {
		return base64.StdEncoding.EncodeToString(r.Payload)
	}
	return string(r.Payload)
}
