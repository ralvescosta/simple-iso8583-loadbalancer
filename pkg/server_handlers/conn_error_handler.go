package server

import "github.com/sirupsen/logrus"

func ConnectionErrorHandler(err error) {
	if err != nil {
		logrus.WithError(err).Error(err.Error())
	}
}
