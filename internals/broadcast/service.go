package broadcast

import (
	"context"

	"github.com/moov-io/iso8583"
	connection "github.com/moov-io/iso8583-connection"
	"github.com/sirupsen/logrus"
)

type (
	BroadcastService interface {
		SetBrandConnection(conn *connection.Connection)
		AddServerConnection(conn *connection.Connection)
		Delivery(ctx context.Context, message *iso8583.Message)
	}

	broadcastService struct {
		brandConnection   *connection.Connection
		serverConnections []*connection.Connection
	}
)

func NewBroadcastService() *broadcastService {
	return &broadcastService{}
}

func (s *broadcastService) SetBrandConnection(conn *connection.Connection) {
	s.brandConnection = conn
}

func (s *broadcastService) AddServerConnection(conn *connection.Connection) {
	s.serverConnections = append(s.serverConnections, conn)
}

func (s *broadcastService) Delivery(ctx context.Context, message *iso8583.Message) {
	go func() {
		for _, serverConn := range s.serverConnections {
			resp, err := serverConn.Send(message)
			if err != nil {
				logrus.
					WithError(err).
					WithField("addr", serverConn.Addr()).
					Error("failed to delivery brand message to this connection")
				continue
			}

			s.brandConnection.Reply(resp)
		}
	}()
}
