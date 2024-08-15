package cmd

import (
	"github.com/moov-io/iso8583"
	"github.com/ralvescosta/simple-iso8583-loadbalancer/internals/broadcast"
	tcpServer "github.com/ralvescosta/simple-iso8583-loadbalancer/pkg/tcp_server"
	"github.com/sirupsen/logrus"
)

func StartISO8583TCPServer(iso8583Spec *iso8583.MessageSpec, broadcastService broadcast.BroadcastService) tcpServer.TCPServer {
	server := tcpServer.NewTCPServer(iso8583Spec, broadcastService)

	go func() {
		logrus.Info("starting iso8583 tcp server")

		if err := server.Start(); err != nil {
			logrus.WithError(err).Fatal("tcp server failed")
		}
	}()

	return server
}
