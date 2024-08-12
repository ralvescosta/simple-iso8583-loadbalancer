package brand

import (
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

func NewBrandTCPConn() (net.Conn, error) {
	brandAddr := os.Getenv("BRAND_ADDR")
	if brandAddr == "" {
		return nil, fmt.Errorf("BRAND_ADDR env must be fulfilled")
	}

	conn, err := net.DialTimeout("tcp", brandAddr, 10)
	if err != nil {
		logrus.WithError(err).Error("failed to connect to brand")
		return nil, err
	}

	return conn, nil
}
