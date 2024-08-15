package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/ralvescosta/simple-iso8583-loadbalancer/cmd"
	"github.com/ralvescosta/simple-iso8583-loadbalancer/internals"
	"github.com/ralvescosta/simple-iso8583-loadbalancer/internals/broadcast"
	brandClient "github.com/ralvescosta/simple-iso8583-loadbalancer/pkg/brand_client"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Fatal("error loading .env file")
	}

	iso8583Spec := internals.NewISO8583Spec()
	broadcastService := broadcast.NewBroadcastService()
	tcpClient, err := brandClient.NewTCPClient(iso8583Spec, broadcastService)
	handleError(err)
	server := cmd.StartISO8583TCPServer(iso8583Spec, broadcastService)
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
