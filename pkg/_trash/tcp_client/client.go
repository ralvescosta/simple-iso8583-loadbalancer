package client

import (
	"net"
	"os"
	"strconv"

	"github.com/moov-io/iso8583"
	proxy "github.com/ralvescosta/simple-iso8583-loadbalancer/internals/proxy_sync"
	"github.com/sirupsen/logrus"
)

type (
	TCPClient interface {
		Start() error
		Send(message *iso8583.Message, respChannel chan *iso8583.Message) error
		Close() error
	}

	tcpClient struct {
		brandConn net.Conn
		proxySync proxy.ProxyMessageSynchronizer

		brandMsgBufferSize int
	}
)

func NewTCPClient(conn net.Conn, proxyMessageSync proxy.ProxyMessageSynchronizer) *tcpClient {
	bufferSize := os.Getenv("BRAND_MSG_BUFFER_SIZE")
	if bufferSize == "" {
		bufferSize = "1024"
	}

	brandMsgBufferSize, _ := strconv.Atoi(bufferSize)

	return &tcpClient{conn, proxyMessageSync, brandMsgBufferSize}
}

func (client *tcpClient) Start() error {
	for {
		readBuffer := make([]byte, client.brandMsgBufferSize)
		lenRead, err := client.brandConn.Read(readBuffer)
		if err != nil {
			logrus.WithError(err).Error("error whiling reading tcp buffer")
			continue
		}

		if err := client.proxySync.Sync(readBuffer[:lenRead]); err != nil {
			logrus.WithError(err).Warn("failed to sync msg")
			continue
		}
	}
}

func (client *tcpClient) Send(message *iso8583.Message, respChannel chan *iso8583.Message) error {
	client.proxySync.AddSync(message, respChannel)

	packedMsg, err := message.Pack()
	if err != nil {
		logrus.WithError(err).Error("failed to packed msg")
		return err
	}

	_, err = client.brandConn.Write(packedMsg)
	if err != nil {
		logrus.WithError(err).Error("failed to write msg brand connection")
		return err
	}

	return nil
}

func (client *tcpClient) Close() error {
	return client.brandConn.Close()
}
