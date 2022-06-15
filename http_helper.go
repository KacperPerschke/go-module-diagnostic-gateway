package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func verifySumDBProxy(host string) bool {
	uri := fmt.Sprintf("%s/supported", host)
	resp, err := callWorld(uri)
	if err != nil {
		globalLogger.WithFields(logrus.Fields{
			`tested_sumdb_host`: host,
		}).Info(err)
		return false
	}
	if resp.StatusCode == 410 {
		globalLogger.WithFields(logrus.Fields{
			`tested_sumdb_host`: host,
		}).Info(`Server returns 410 == Gone`)
		return false
	}
	return true
}

func verifyModuleProxy(host string) bool {
	uri := fmt.Sprintf("%s/rsc.io/quote/@v/list", host)
	resp, err := callWorld(uri)
	if err != nil {
		globalLogger.WithFields(logrus.Fields{
			`tested_sumdb_host`: host,
		}).Info(err)
		return false
	}
	if resp.StatusCode != 200 {
		globalLogger.WithFields(logrus.Fields{
			"tested_sumdb_host":    host,
			"http.response.body":   resp.toLog(false),
			"http.response.status": resp.StatusCode,
		}).Info(`Server returns not 200.`)
		return false
	}
	return true

}
