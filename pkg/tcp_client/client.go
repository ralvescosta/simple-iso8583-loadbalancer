package anewway

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/moov-io/iso8583"
	connection "github.com/moov-io/iso8583-connection"
	"github.com/moov-io/iso8583/network"
	"github.com/sirupsen/logrus"
)

type (
	TCPClient interface {
		Start() error
		Send(ctx context.Context, message *iso8583.Message) (*iso8583.Message, error)
		Close() error
	}

	tcpClient struct {
		spec                   *iso8583.MessageSpec
		brandAddr              string
		brandReconnectWait     time.Duration
		brandConnectionTimeout time.Duration

		pool *connection.Pool
	}
)

func NewTCPClient(spec *iso8583.MessageSpec) (*tcpClient, error) {
	c := tcpClient{}

	if err := c.loadEnvs(); err != nil {
		return nil, err
	}

	pool, err := connection.NewPool(
		c.factory,
		[]string{c.brandAddr},
		connection.PoolMinConnections(1),
		connection.PoolReconnectWait(c.brandReconnectWait*time.Second),
		connection.PoolErrorHandler(c.connectionErrorHandler),
	)

	if err != nil {
		return nil, fmt.Errorf("creating iso8583 connection: %w", err)
	}

	c.pool = pool

	return &c, nil
}

func (client *tcpClient) Start() error {
	return client.pool.Connect()
}

func (client *tcpClient) Close() error {
	return client.pool.Close()
}

func (client *tcpClient) Send(ctx context.Context, message *iso8583.Message) (*iso8583.Message, error) {
	conn, err := client.pool.Get()
	if err != nil {
		logrus.WithError(err).Error("failed to retrieve connection from pool")
		return nil, err
	}

	resp, err := conn.Send(message)
	if err != nil {
		logrus.WithError(err).Error("failed to send msg to brand")
		return nil, err
	}

	return resp, nil
}

func (client *tcpClient) factory(addr string) (*connection.Connection, error) {
	c, err := connection.New(
		addr,
		client.spec,
		client.readMessageLength,
		client.writeMessageLength,
		connection.ConnectTimeout(client.brandConnectionTimeout*time.Second),
		connection.ReadTimeout(10*time.Second),
		connection.ErrorHandler(client.connectionErrorHandler),
		connection.ConnectionClosedHandler(client.connectionClosedHandler),
		connection.ConnectionEstablishedHandler(client.connectionEstablishedHandler),
		connection.InboundMessageHandler(client.inboundMessageHandler),
	)
	if err != nil {
		return nil, fmt.Errorf("building iso8583 connection: %w", err)
	}

	return c, nil

}

func (s *tcpClient) readMessageLength(r io.Reader) (int, error) {
	header := network.NewBinary2BytesHeader()
	n, err := header.ReadFrom(r)
	if err != nil {
		return n, err
	}

	return header.Length(), nil
}

func (s *tcpClient) writeMessageLength(w io.Writer, length int) (int, error) {
	header := network.NewBinary2BytesHeader()
	header.SetLength(length)
	return header.WriteTo(w)
}

func (s *tcpClient) connectionEstablishedHandler(c *connection.Connection) {
	logrus.WithField("host", c.Addr()).Info("client connected")
}

func (s *tcpClient) inboundMessageHandler(c *connection.Connection, message *iso8583.Message) {
	_ = context.Background()
}

func (s *tcpClient) connectionErrorHandler(err error) {
	if err != nil {
		logrus.WithError(err).Error(err.Error())
	}
}

func (s *tcpClient) connectionClosedHandler(c *connection.Connection) {}

func (c *tcpClient) loadEnvs() error {
	brandAddr := os.Getenv("BRAND_ADDR")
	if brandAddr == "" {
		return fmt.Errorf("BRAND_ADDR env must be fulfilled")
	}

	//
	reconnectWait := os.Getenv("BRAND_RECONNECT_WAIT")
	if reconnectWait == "" {
		reconnectWait = "10"
	}

	brandReconnectWait, _ := strconv.Atoi(reconnectWait)

	//
	connectionTimeout := os.Getenv("BRAND_CONNECTION_TIMEOUT")
	if connectionTimeout == "" {
		connectionTimeout = "10"
	}

	brandConnectionTimeout, _ := strconv.Atoi(connectionTimeout)

	//
	c.brandAddr = brandAddr
	c.brandReconnectWait = time.Duration(brandReconnectWait)
	c.brandConnectionTimeout = time.Duration(brandConnectionTimeout)

	return nil
}
