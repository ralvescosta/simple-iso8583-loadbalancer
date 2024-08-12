package cmd

import (
	"fmt"
	"os"

	"github.com/moov-io/iso8583"
	connection "github.com/moov-io/iso8583-connection"
	server "github.com/moov-io/iso8583-connection/server"
	proxy "github.com/ralvescosta/mastercard-tcp-proxy/internals/proxy_handler"
	tcpHandlers "github.com/ralvescosta/mastercard-tcp-proxy/pkg/server_handlers"
	"github.com/sirupsen/logrus"
)

func StartISO8583TCPServer(iso8583Spec *iso8583.MessageSpec, proxyHandler proxy.ProxyMessageHandler) *server.Server {
	srv := server.New(
		iso8583Spec,
		tcpHandlers.ReadMessageLength,
		tcpHandlers.WriteMessageLength,
		connection.ConnectionEstablishedHandler(tcpHandlers.ConnectionEstablishedHandler),
		connection.InboundMessageHandler(tcpHandlers.InboundMessageHandler(proxyHandler)),
		connection.ErrorHandler(tcpHandlers.ConnectionErrorHandler),
		connection.ConnectionClosedHandler(tcpHandlers.ConnectionClosedHandler),
	)

	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4001"
	}

	addr := fmt.Sprintf("%v:%v", host, port)

	go func() {
		logrus.Info("starting iso8583 tcp server")

		if err := srv.Start(addr); err != nil {
			logrus.WithError(err).Fatal("tcp server failed")
		}
	}()

	return srv
}
