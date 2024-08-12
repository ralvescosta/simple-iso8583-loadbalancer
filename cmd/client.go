package cmd

import (
	tcpClient "github.com/ralvescosta/mastercard-tcp-proxy/pkg/tcp_client"
	"github.com/sirupsen/logrus"
)

func StartTCPClient(client tcpClient.TCPClient) tcpClient.TCPClient {
	go func() {
		if err := client.Start(); err != nil {
			logrus.WithError(err).Fatal("brand connection failed")
		}
	}()

	return client
}
