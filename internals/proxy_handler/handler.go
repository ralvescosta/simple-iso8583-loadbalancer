package proxy

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/moov-io/iso8583"
	connection "github.com/moov-io/iso8583-connection"
	tcpClient "github.com/ralvescosta/simple-iso8583-loadbalancer/pkg/tcp_client"
)

type (
	ProxyMessageHandler interface {
		Handler(ctx context.Context, c *connection.Connection, message *iso8583.Message)
	}

	proxyMessageHandler struct {
		client      tcpClient.TCPClient
		respTimeout time.Duration
	}
)

func NewProxyMessageHandler(tcpClient tcpClient.TCPClient) *proxyMessageHandler {
	brandRespTimeout := os.Getenv("BRAND_RESPONSE_TIMEOUT")
	if brandRespTimeout == "" {
		brandRespTimeout = "10"
	}

	respTimeout, _ := strconv.Atoi(brandRespTimeout)

	return &proxyMessageHandler{
		client:      tcpClient,
		respTimeout: time.Duration(respTimeout) * time.Second,
	}
}

func (handler *proxyMessageHandler) Handler(ctx context.Context, c *connection.Connection, message *iso8583.Message) {
	// respChannel := make(chan *iso8583.Message)
	// maxLimitCtx, cancelFunc := context.WithTimeout(ctx, handler.respTimeout)
	// defer cancelFunc()

	// if err := handler.client.Send(message, respChannel); err != nil {
	// 	logrus.WithError(err).Error("failed to delivery msg to brand")
	// 	return
	// }

	// resp := &iso8583.Message{}
	// select {
	// case <-maxLimitCtx.Done():
	// 	logrus.
	// 		WithContext(ctx).
	// 		WithError(maxLimitCtx.Err()).
	// 		Error("brand did not respond this transaction")
	// case resp = <-respChannel:
	// }

	// if _, err := c.Send(resp, connection.SendTimeout(5*time.Second)); err != nil {
	// 	logrus.
	// 		WithContext(ctx).
	// 		WithError(err).
	// 		WithField("clientAddr", c.Addr()).
	// 		Error("failed to delivery the brand response")
	// }
}
