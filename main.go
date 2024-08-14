package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/ralvescosta/simple-iso8583-loadbalancer/cmd"
	"github.com/ralvescosta/simple-iso8583-loadbalancer/internals"
	proxyMessageHandler "github.com/ralvescosta/simple-iso8583-loadbalancer/internals/proxy_handler"
	tcpClient "github.com/ralvescosta/simple-iso8583-loadbalancer/pkg/tcp_client"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Fatal("error loading .env file")
	}

	iso8583Spec := internals.NewISO8583Spec()
	tcpClient, err := tcpClient.NewTCPClient(iso8583Spec)
	handleError(err)

	handler := proxyMessageHandler.NewProxyMessageHandler(tcpClient)
	server := cmd.StartISO8583TCPServer(iso8583Spec, handler)
	client := cmd.StartTCPClient(tcpClient)

	shotdown := make(chan os.Signal, 1)
	signal.Notify(shotdown, syscall.SIGINT, syscall.SIGTERM)

	<-shotdown

	logrus.Warn("starting shotdown...")

	logrus.Warn("closing server...")
	server.Close()

	time.Sleep(2 * time.Second)

	logrus.Warn("closing brand connection...")
	client.Close()

	logrus.Info("shotdown finished!")
}

func handleError(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}
