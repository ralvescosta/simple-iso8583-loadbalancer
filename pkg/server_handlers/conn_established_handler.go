package server

import (
	connection "github.com/moov-io/iso8583-connection"
	"github.com/sirupsen/logrus"
)

func ConnectionEstablishedHandler(c *connection.Connection) {
	logrus.WithField("host", c.Addr()).Info("client connected")
}
