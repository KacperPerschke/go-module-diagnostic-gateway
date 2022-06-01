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
	return urlsUnverified
}

func listSumDBProxies(confReadIn *viper.Viper) listOfHosts {
	urlsUnverified := confReadIn.GetStringSlice(`sumdb_proxies`)
	useModProxsAlso := confReadIn.GetBool(`use_module_proxies_as_sumdb_also`)
	if useModProxsAlso {
		urlsToModProxs := confReadIn.GetStringSlice(`module_proxies`)
		urlsUnverified = append(urlsUnverified, urlsToModProxs...)
	}
	return urlsUnverified
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
