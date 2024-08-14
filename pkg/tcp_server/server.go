package server

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/moov-io/iso8583"
	connection "github.com/moov-io/iso8583-connection"
	iso8583Server "github.com/moov-io/iso8583-connection/server"
	"github.com/moov-io/iso8583/network"
	proxy "github.com/ralvescosta/simple-iso8583-loadbalancer/internals/proxy_handler"
	tcpClient "github.com/ralvescosta/simple-iso8583-loadbalancer/pkg/tcp_client"
	"github.com/sirupsen/logrus"
)

type (
	TCPServer interface {
		Start() error
		Close()
	}

	tcpServer struct {
		iso8583Spec   *iso8583.MessageSpec
		proxyHandler  proxy.ProxyMessageHandler
		serverAddr    string
		iso8583Server *iso8583Server.Server
		tcpClient     tcpClient.TCPClient
	}
)

func NewTCPServer(iso8583Spec *iso8583.MessageSpec, proxyHandler proxy.ProxyMessageHandler) *tcpServer {
	host := os.Getenv("TCP_SERVER_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("TCP_SERVER_PORT")
	if port == "" {
		port = "4001"
	}

	srv := &tcpServer{
		iso8583Spec:  iso8583Spec,
		proxyHandler: proxyHandler,
		serverAddr:   fmt.Sprintf("%v:%v", host, port),
	}

	srv.createISO8583ServerInstance()

	return srv
}

func (s *tcpServer) Start() error {
	return s.iso8583Server.Start(s.serverAddr)
}

func (s *tcpServer) Close() {
	s.iso8583Server.Close()
}

func (s *tcpServer) createISO8583ServerInstance() {
	s.iso8583Server = iso8583Server.New(
		s.iso8583Spec,
		s.readMessageLength,
		s.writeMessageLength,
		connection.ConnectionEstablishedHandler(s.connectionEstablishedHandler),
		connection.InboundMessageHandler(s.inboundMessageHandler),
		connection.ErrorHandler(s.connectionErrorHandler),
		connection.ConnectionClosedHandler(s.connectionClosedHandler),
	)
}

func (s *tcpServer) inboundMessageHandler(c *connection.Connection, message *iso8583.Message) {
	ctx := context.Background()
	resp, err := s.tcpClient.Send(ctx, message)
	if err != nil {
		logrus.WithError(err).Error("failed to send msg to brand")

		c.Reply(resp)

		return
	}

	c.Reply(resp)
}

func (s *tcpServer) readMessageLength(r io.Reader) (int, error) {
	header := network.NewBinary2BytesHeader()
	n, err := header.ReadFrom(r)
	if err != nil {
		return n, err
	}

	return header.Length(), nil
}

func (s *tcpServer) writeMessageLength(w io.Writer, length int) (int, error) {
	header := network.NewBinary2BytesHeader()
	header.SetLength(length)
	return header.WriteTo(w)
}

func (s *tcpServer) connectionEstablishedHandler(c *connection.Connection) {
	logrus.WithField("host", c.Addr()).Info("client connected")
}

func (s *tcpServer) connectionErrorHandler(err error) {
	if err != nil {
		logrus.WithError(err).Error(err.Error())
	}
}

func (s *tcpServer) connectionClosedHandler(c *connection.Connection) {}
