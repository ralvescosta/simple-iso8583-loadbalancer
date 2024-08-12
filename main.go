package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/ralvescosta/mastercard-tcp-proxy/cmd"
	"github.com/ralvescosta/mastercard-tcp-proxy/internals"
	proxyMessageHandler "github.com/ralvescosta/mastercard-tcp-proxy/internals/proxy_handler"
	proxySynchronizer "github.com/ralvescosta/mastercard-tcp-proxy/internals/proxy_sync"
	brandConnector "github.com/ralvescosta/mastercard-tcp-proxy/pkg/brand_conn"
	tcpClient "github.com/ralvescosta/mastercard-tcp-proxy/pkg/tcp_client"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Fatal("error loading .env file")
	}

	brandConn, err := brandConnector.NewBrandTCPConn()
	if err != nil {
		logrus.Fatal(err)
	}

	iso8583Spec := internals.NewISO8583Spec()
	synchronizer := proxySynchronizer.NewProxyMessageSynchronizer(iso8583Spec)
	tcpClient := tcpClient.NewTCPClient(brandConn, synchronizer)
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
