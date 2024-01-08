/*
Copyright © 2024 Alexander Sharov <kvendingoldo@gmail.com>, Nikolai Mishin <sanduku.default@gmail.com>, Anastasiia Kozlova <anastasiia.kozlova245@gmail.com>
*/
package main

import (
	"github.com/opentofuutils/tenv/cmd"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetReportCaller(false)

	var formatter log.Formatter

	if formatterType, ok := os.LookupEnv("FORMATTER_TYPE"); ok {
		if formatterType == "JSON" {
			formatter = &log.JSONFormatter{PrettyPrint: false}
		}

		if formatterType == "TEXT" {
			formatter = &log.TextFormatter{DisableColors: false}
		}
	}

	if formatter == nil {
		formatter = &log.TextFormatter{DisableColors: false}
	}

	log.SetFormatter(formatter)

	var logLevel log.Level
	var err error

	if ll, ok := os.LookupEnv("LOG_LEVEL"); ok {
		logLevel, err = log.ParseLevel(ll)
		if err != nil {
			logLevel = log.DebugLevel
		}
	} else {
		logLevel = log.DebugLevel
	}

	log.SetLevel(logLevel)
}

func main() {
	cmd.Execute()
}
