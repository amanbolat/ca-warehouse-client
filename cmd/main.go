package main

import (
	"github.com/amanbolat/ca-warehouse-client/config"
	"github.com/amanbolat/ca-warehouse-client/server"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

var conf config.Config
var logger *logrus.Logger

func main() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	if conf.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	err := envconfig.Process("", &conf)
	if err != nil {
		logger.Fatalf("could not parse env vars: %v", err)
	}

	s := server.NewServer(conf)

	port := 11201
	if conf.Port > 1024 {
		port = 11201
	}
	s.Start(port)
}
