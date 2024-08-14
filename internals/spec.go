package internals

import (
	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
)

func NewISO8583Spec() *iso8583.MessageSpec {
	return &iso8583.MessageSpec{
		Fields: map[int]field.Field{
			0: field.NewString(&field.Spec{
				Length:      4,
				Description: "Message Type Indicator",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			1: field.NewBitmap(&field.Spec{
				Description: "Bitmap",
				Enc:         encoding.BytesToASCIIHex,
				Pref:        prefix.Hex.Fixed,
			}),

			// Message fields:
			2: field.NewString(&field.Spec{
				Length:      19,
				Description: "Primary Account Number",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			3: field.NewNumeric(&field.Spec{
				Length:      6,
				Description: "Processing Code",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			7: field.NewString(&field.Spec{
				Length:      10,
				Description: "Transmission date & time",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			11: field.NewString(&field.Spec{
				Length:      6,
				Description: "System trace audit number (STAN)",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			12: field.NewString(&field.Spec{
				Length:      6,
				Description: "Local transaction time (hhmmss)",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			13: field.NewString(&field.Spec{
				Length:      4,
				Description: "Local Transaction Date (MMDD)",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
		},
	}
}
