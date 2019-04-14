package main

import (
	"flag"
	"net/http"

	"github.com/leominov/vault-auth-proxy/pkg/sso"
	"github.com/sirupsen/logrus"
)

var (
	configFile    = flag.String("config", "config.yaml", "Path to configuration file")
	listenAddress = flag.String("listen-address", ":8080", "Address to server requests")
	logLevel      = flag.String("log-level", "debug", "Logging level")
)

func main() {
	flag.Parse()

	logger := logrus.New()
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logger.Fatal(err)
	}

	logger.SetLevel(level)
	logger.Info("Starting vault-auth-proxy...")

	c, err := sso.LoadConfig(*configFile)
	if err != nil {
		logger.Fatal(err)
	}

	sso, err := sso.New(c, logger)
	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Addr:    *listenAddress,
		Handler: sso,
	}

	logger.Infof("Listening address: %s", server.Addr)
	logger.Fatal(server.ListenAndServe())
}