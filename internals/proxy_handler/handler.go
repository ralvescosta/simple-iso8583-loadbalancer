package proxy

import (
	"context"
	"time"

	"github.com/moov-io/iso8583"
	connection "github.com/moov-io/iso8583-connection"
	tcpClient "github.com/ralvescosta/mastercard-tcp-proxy/pkg/tcp_client"
	"github.com/sirupsen/logrus"
)

type (
	ProxyMessageHandler interface {
		Handler(ctx context.Context, c *connection.Connection, message *iso8583.Message)
	}

	proxyMessageHandler struct {
		client tcpClient.TCPClient
	}
)

func NewProxyMessageHandler(tcpClient tcpClient.TCPClient) *proxyMessageHandler {
	return &proxyMessageHandler{tcpClient}
}

func (handler *proxyMessageHandler) Handler(ctx context.Context, c *connection.Connection, message *iso8583.Message) {
	respChannel := make(chan *iso8583.Message)
	maxLimitCtx, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFunc()

	if err := handler.client.Send(message, respChannel); err != nil {
		logrus.WithError(err).Error("failed to delivery msg to brand")
		return
	}

	resp := &iso8583.Message{}
	select {
	case <-maxLimitCtx.Done():
		logrus.
			WithContext(ctx).
			WithError(maxLimitCtx.Err()).
			Error("brand did not respond this transaction")
	case resp = <-respChannel:
	}

	if _, err := c.Send(resp, connection.SendTimeout(5*time.Second)); err != nil {
		logrus.
			WithContext(ctx).
			WithError(err).
			WithField("clientAddr", c.Addr()).
			Error("failed to delivery the brand response")
	}
}
