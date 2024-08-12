package server

import (
	"io"

	"github.com/moov-io/iso8583/network"
)

func ReadMessageLength(r io.Reader) (int, error) {
	header := network.NewBinary2BytesHeader()
	n, err := header.ReadFrom(r)
	if err != nil {
		return n, err
	}

	return header.Length(), nil
}
