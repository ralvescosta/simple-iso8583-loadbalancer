package cmd

import (
	"github.com/moov-io/iso8583"
	proxy "github.com/ralvescosta/simple-iso8583-loadbalancer/internals/proxy_handler"
	tcpServer "github.com/ralvescosta/simple-iso8583-loadbalancer/pkg/tcp_server"
	"github.com/sirupsen/logrus"
)

func StartISO8583TCPServer(iso8583Spec *iso8583.MessageSpec, proxyHandler proxy.ProxyMessageHandler) tcpServer.TCPServer {
	server := tcpServer.NewTCPServer(iso8583Spec, proxyHandler)

	go func() {
		logrus.Info("starting iso8583 tcp server")

		if err := server.Start(); err != nil {
			logrus.WithError(err).Fatal("tcp server failed")
		}
	}()

	return server
}
