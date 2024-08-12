package server

import (
	"io"

	"github.com/moov-io/iso8583/network"
)

func WriteMessageLength(w io.Writer, length int) (int, error) {
	header := network.NewBinary2BytesHeader()
	header.SetLength(length)
	return header.WriteTo(w)
}
