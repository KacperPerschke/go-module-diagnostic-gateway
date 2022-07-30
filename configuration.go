package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type configInterface interface {
	getModuleProxies() listOfHosts
	getSumDBProxies() listOfHosts
}

type configStructure struct {
	moduleProxies listOfHosts
	sumDBProxies  listOfHosts
}

type listOfHosts []string

var globalConfig configInterface

func (c configStructure) getModuleProxies() listOfHosts {
	return c.moduleProxies
}

func (c configStructure) getSumDBProxies() listOfHosts {
	return c.sumDBProxies
}

func listModuleProxies(confReadIn *viper.Viper) listOfHosts {
	urlsUnverified := confReadIn.GetStringSlice(`module_proxies`)
	globalLogger.WithFields(logrus.Fields{
		`list of proposed module proxies`: urlsUnverified,
	}).Info("Parsed Module Proxies configuration")
	urlsVerified, err := verifyListOfHosts(
		urlsUnverified,
		verifyModuleProxy,
	)
	if err != nil {
		globalLogger.Fatal(err)
	}
	globalLogger.WithFields(logrus.Fields{
		`list of verified sumdb proxies`: urlsVerified,
	}).Info("Parsed SumDB Proxies configuration")
	return urlsVerified
}

func listSumDBProxies(confReadIn *viper.Viper) listOfHosts {
	urlsUnverified := confReadIn.GetStringSlice(`sumdb_proxies`)
	if confReadIn.GetBool(`use_module_proxies_as_sumdb_also`) {
		urlsToModProxs := confReadIn.GetStringSlice(`module_proxies`)
		urlsUnverified = append(urlsUnverified, urlsToModProxs...)
	}
	globalLogger.WithFields(logrus.Fields{
		`list of proposed sumdb proxies`: urlsUnverified,
	}).Info("Parsed SumDB Proxies configuration")
	urlsVerified, err := verifyListOfHosts(
		urlsUnverified,
		verifySumDBProxy,
	)
	if err != nil {
		globalLogger.Fatal(err)
	}
	globalLogger.WithFields(logrus.Fields{
		`list of verified sumdb proxies`: urlsVerified,
	}).Info("Parsed SumDB Proxies configuration")
	return urlsVerified
}

func produceConfiguration(confFilePath string) configInterface {
	configSlurped := slurpConfiguration(confFilePath)

	globalLogger.Info("Starts parsing configuration.")
	configPrepared := configStructure{
		moduleProxies: listModuleProxies(configSlurped),
		sumDBProxies:  listSumDBProxies(configSlurped),
	}
	globalLogger.Info("Configuration parsed successfully.")
	return &configPrepared // → https://play.golang.org/p/kIhvK8m62aG
}

func slurpConfiguration(confFilePath string) *viper.Viper {
	globalLogger.WithFields(logrus.Fields{
		"config file path": confFilePath,
	}).Info("Starts reading configuration from file.")
	confParser := viper.New()
	confParser.AddConfigPath(".")
	confParser.SetConfigName(confFilePath)
	confParser.SetConfigType("yaml")
	err := confParser.ReadInConfig()
	if err != nil {
		globalLogger.Fatal(err)
	}
	globalLogger.Info("Configuration read in successfully.")
	return confParser
}
