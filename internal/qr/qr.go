package qr

import (
	"fmt"
	"image/png"
	"io"
	"io/ioutil"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// Encode converts the content of io.Reader into QR Code and writes it to the io.Writer
func Encode(r io.Reader, w io.Writer) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("could not read io.Reader content : %w", err)
	}
	return EncodeString(string(data), w)
}

// EncodeString converts the string content into QR Code and writes it to the io.Writer as PNG
func EncodeString(content string, w io.Writer) error {
	qrc, err := qr.Encode(content, qr.M, qr.Auto)
	qrc, err = barcode.Scale(qrc, 250, 250)
	if err != nil {
		return fmt.Errorf("could not generate QRCode: %w", err)
	}
	return png.Encode(w, qrc)
}
