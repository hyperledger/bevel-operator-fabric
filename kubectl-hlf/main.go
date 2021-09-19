package main

import (
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "info"
	}
	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}
	// set global log level
	logrus.SetLevel(ll)
	if err := cmd.NewCmdHLF().Execute(); err != nil {
		os.Exit(1)
	}
}
