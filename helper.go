package main

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

type hostVerificator func(string) bool

func verifyListOfHosts(unverified listOfHosts, verificator hostVerificator) (listOfHosts, error) {
	verified := listOfHosts{}
	verificationError := *new(error)
	for _, host := range unverified {
		if verificator(host) {
			verified = append(verified, host)
		}
	}
	if len(verified) == 0 {
		return listOfHosts{}, errors.New("none of the hosts passed verification")
	}
	return verified, verificationError
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
