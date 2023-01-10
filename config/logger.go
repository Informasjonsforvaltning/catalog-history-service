package config

import "github.com/sirupsen/logrus"

func LoggerSetup() {
	formatter := logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "severity",
		},
	}
	logrus.SetFormatter(&formatter)
}
