package cmd

import (
	brandClient "github.com/ralvescosta/simple-iso8583-loadbalancer/pkg/brand_client"
	"github.com/sirupsen/logrus"
)

func StartTCPClient(client brandClient.TCPClient) brandClient.TCPClient {
	go func() {
		if err := client.Start(); err != nil {
			logrus.WithError(err).Fatal("brand connection failed")
		}
	}()

	return client
}
