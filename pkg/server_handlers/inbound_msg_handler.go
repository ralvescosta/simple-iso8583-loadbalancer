package server

import (
	"context"

	"github.com/moov-io/iso8583"
	connection "github.com/moov-io/iso8583-connection"
	proxy "github.com/ralvescosta/mastercard-tcp-proxy/internals/proxy_handler"
)

func InboundMessageHandler(handler proxy.ProxyMessageHandler) func(c *connection.Connection, message *iso8583.Message) {
	return func(c *connection.Connection, message *iso8583.Message) {
		ctx := context.Background()

		handler.Handler(ctx, c, message)
	}
}
