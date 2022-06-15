package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
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

func callWorld(uri string) (myRespType, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return myRespType{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	myResp := myRespType{
		Header:     resp.Header.Clone(),
		Payload:    b,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
	}
	return myResp, err
}

func (r myRespType) copyHeadersTo(w http.ResponseWriter) {
	for n, values := range r.Header {
		for _, v := range values {
			w.Header().Add(n, v)
		}
	}
}

func helperForModuleProtocol(host, path string) (myRespType, error) {
	uri := fmt.Sprintf("%s%s", host, path)
	resp, err := callWorld(uri)
	if err != nil {
		globalLogger.WithFields(logrus.Fields{
			"http.request.uri":    uri,
			"http.response.error": err.Error(),
		}).Error(`Didn't get proper answear from external server`)
		return resp, nil
	}
	respContentType := resp.Header.Get("Content-Type")
	logAsBase64 := bool(respContentType == "application/zip")
	globalLogger.WithFields(logrus.Fields{
		"http.request.uri":           uri,
		"http.response.body":         resp.toLog(logAsBase64),
		"http.response.status":       resp.StatusCode,
		"http.response.Content-Type": respContentType,
	}).Debug(`Got proper answear.`)
	return resp, nil
}

func (r myRespType) toLog(logAsBase64 bool) string {
	if logAsBase64 {
		return base64.StdEncoding.EncodeToString(r.Payload)
	}
	return string(r.Payload)
}
