package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/moov-io/iso8583"
	connection "github.com/moov-io/iso8583-connection"
	iso8583Server "github.com/moov-io/iso8583-connection/server"
	"github.com/moov-io/iso8583/network"
	"github.com/ralvescosta/simple-iso8583-loadbalancer/internals/broadcast"
	brandClient "github.com/ralvescosta/simple-iso8583-loadbalancer/pkg/brand_client"
	"github.com/sirupsen/logrus"
)

type (
	TCPServer interface {
		Start() error
		Close()
	}

	tcpServer struct {
		iso8583Spec      *iso8583.MessageSpec
		serverAddr       string
		iso8583Server    *iso8583Server.Server
		brandClient      brandClient.TCPClient
		broadcastService broadcast.BroadcastService
	}
)

func NewTCPServer(iso8583Spec *iso8583.MessageSpec, broadcastService broadcast.BroadcastService) *tcpServer {
	host := os.Getenv("TCP_SERVER_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("TCP_SERVER_PORT")
	if port == "" {
		port = "4001"
	}

	srv := &tcpServer{
		iso8583Spec: iso8583Spec,
		serverAddr:  fmt.Sprintf("%v:%v", host, port),
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
		//@WARNING github.com/moov-io/iso8583-connection ConnectionEstablishedHandler are not working
		//instead we need to use github.com/moov-io/iso8583-connection/server.Server.AddConnectionHandler(...)
		connection.ConnectionEstablishedHandler(s.connectionEstablishedHandler),
		connection.InboundMessageHandler(s.inboundMessageHandler),
		connection.ErrorHandler(s.connectionErrorHandler),
		connection.ConnectionClosedHandler(s.connectionClosedHandler),
	)

	s.iso8583Server.AddConnectionHandler(func(tcpConn net.Conn) {
		conn, err := connection.NewFrom(tcpConn, s.iso8583Spec, s.readMessageLength, s.writeMessageLength)
		if err != nil {
			logrus.
				WithError(err).
				Error("failed to create the abstracted moov-io connection")
			return
		}

		s.broadcastService.AddServerConnection(conn)
	})
}

func (s *tcpServer) inboundMessageHandler(c *connection.Connection, message *iso8583.Message) {
	ctx := context.Background()

	resp, err := s.brandClient.Send(ctx, message)
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
	err := header.SetLength(length)
	if err != nil {
		return 0, err
	}

	_, err = header.WriteTo(w)
	if err != nil {
		return 0, err
	}

	return header.Length(), nil
}

func (s *tcpServer) connectionEstablishedHandler(c *connection.Connection) {
	logrus.WithField("host", c.Addr()).Info("client connected")
	s.broadcastService.AddServerConnection(c)
}

func (s *tcpServer) connectionErrorHandler(err error) {
	if err != nil {
		logrus.WithError(err).Error(err.Error())
	}
}

func (s *tcpServer) connectionClosedHandler(c *connection.Connection) {
	logrus.WithField("host", c.Addr()).Info("client closed")
}
