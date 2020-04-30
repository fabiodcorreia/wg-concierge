package qr

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/yeqown/go-qrcode"
)

// Encode converts the content of io.Reader into QR Code and writes it to the io.Writer
func Encode(r io.Reader, w io.Writer) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("could not read io.Reader content : %w", err)
	}
	return EncodeString(string(data), w)
}

// EncodeString converts the string content into QR Code and writes it to the io.Writer
func EncodeString(content string, w io.Writer) error {
	qrc, err := qrcode.New(content)
	if err != nil {
		return fmt.Errorf("could not generate QRCode: %w", err)
	}
	if err := qrc.SaveTo(w); err != nil {
		return fmt.Errorf("could not save image: %w", err)
	}
	return nil
}
