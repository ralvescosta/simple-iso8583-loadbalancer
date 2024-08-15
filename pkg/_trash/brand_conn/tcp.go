package brand

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func NewBrandTCPConn() (net.Conn, error) {
	brandAddr := os.Getenv("BRAND_ADDR")
	if brandAddr == "" {
		return nil, fmt.Errorf("BRAND_ADDR env must be fulfilled")
	}

	connectionTimeout := os.Getenv("BRAND_CONNECTION_TIMEOUT")
	if connectionTimeout == "" {
		connectionTimeout = "10"
	}

	brandConnectionTimeout, _ := strconv.Atoi(connectionTimeout)

	conn, err := net.DialTimeout("tcp", brandAddr, time.Duration(brandConnectionTimeout)*time.Second)
	if err != nil {
		logrus.WithError(err).Error("failed to connect to brand")
		return nil, err
	}

	return conn, nil
}
