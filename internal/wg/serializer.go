package wg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"text/template"
)

const clientConfig = `
[Interface]
Address = {{ .Interface.Address }}
PrivateKey = {{ .Interface.PrivateKey }}
DNS = {{ .Interface.DNS }}

[Peer]
PublicKey = {{ .Peer.PublicKey }}
Endpoint = {{ .Peer.Endpoint }}
AllowedIPs = {{ .Peer.AllowedIPs }}
PersistentKeepalive = {{ .Peer.KeepAlive }}
`

// MarshalClient will convert a ClientConfig Value in text and serialize it on the io.Writer
func MarshalClient(cc ClientConfig, w io.Writer) error {
	t, err := template.New("clientConfig").Parse(clientConfig)
	if err != nil {
		return fmt.Errorf("fail to compile client configuration template: %w", err)
	}

	err = t.Execute(w, cc)
	if err != nil {
		return fmt.Errorf("fail to marshal client configuration: %w", err)
	}
	return nil
}

// MarshalClientToString will convert a ClientConfig Value in text and return a string with the content
func MarshalClientToString(cc ClientConfig) (result string, err error) {
	var tpl bytes.Buffer
	err = MarshalClient(cc, &tpl)
	if err != nil {
		return result, err // Not need to wrap because the returned error is from a local function already wrapped
	}

	return tpl.String(), nil
}

// UnmarshalClient will convert a Client Configuration from io.Reader into a ClientConfig Value
func UnmarshalClient(r io.Reader, cc *ClientConfig) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return nil
}
