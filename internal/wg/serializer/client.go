package serializer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"text/template"

	"github.com/fabiodcorreia/wg-concierge/internal/wg"
)

//
//
//
//
// Marshal and Unmarshal the Client Configuration File (Imported on the client peer)

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

//! Talvez so precise de serializar os Client Confs (Browser e File) e deserializar os Server Confs para ir buscar a Priv e gerar a Pub Keys e ir buscar os IPs

// MarshalClient will convert a ClientConfig Value in text and serialize it on the io.Writer
func MarshalClient(cc wg.ClientConfig, w io.Writer) error {
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

// MarshalClientToStr will convert a ClientConfig Value in text and return a string with the content
func MarshalClientToStr(cc wg.ClientConfig) (result string, err error) {
	var tpl bytes.Buffer
	err = MarshalClient(cc, &tpl)
	if err != nil {
		return result, err // Not need to wrap because the returned error is from a local function already wrapped
	}

	return tpl.String(), nil
}

// Client Interface Properties
var regexClientAddresss = regexp.MustCompile(`Address\s*=\s*(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\/\d{2})`)
var regexClientPrivKey = regexp.MustCompile(`PrivateKey\s*=\s*(.{44})`)
var regexClientDNS = regexp.MustCompile(`DNS\s*=\s*(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)

// Client Peer Properties
var regexClientPubKey = regexp.MustCompile(`PublicKey\s*=\s*(.{44})`)
var regexClientEndpoint = regexp.MustCompile(`Endpoint\s*=\s*(.+:\d+)`)
var regexClientAllowIP = regexp.MustCompile(`AllowedIPs\s*=\s*(.+)`)

// UnmarshalClient will convert a Client Configuration from io.Reader into a ClientConfig Value
func UnmarshalClient(r io.Reader, cc *wg.ClientConfig) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("fail to unmarshal client configuration: %w", err)
	}

	return UnmarshalClientFromStr(string(buf), cc)
}

// UnmarshalClientFromStr will convert a Client Configuration from io.Reader into a ClientConfig Value
func UnmarshalClientFromStr(c string, cc *wg.ClientConfig) error {
	// Interface
	val, err := findConfig(&c, regexClientAddresss, true)
	if err != nil {
		return fmt.Errorf("client interface address: %w", err)
	}
	cc.Interface.Address = val

	val, err = findConfig(&c, regexClientPrivKey, true)
	if err != nil {
		return fmt.Errorf("client interface private key: %w", err)
	}
	cc.Interface.PrivateKey = val

	val, err = findConfig(&c, regexClientDNS, false)
	if err != nil {
		return fmt.Errorf("client interface dns: %w", err)
	}
	cc.Interface.DNS = val

	// Peer
	val, err = findConfig(&c, regexClientPubKey, true)
	if err != nil {
		return fmt.Errorf("client peer public key: %w", err)
	}
	cc.Peer.PublicKey = val

	val, err = findConfig(&c, regexClientEndpoint, true)
	if err != nil {
		return fmt.Errorf("client peer endpoint: %w", err)
	}
	cc.Peer.Endpoint = val

	val, err = findConfig(&c, regexClientAllowIP, true)
	if err != nil {
		return fmt.Errorf("client peer allowed ips: %w", err)
	}
	cc.Peer.AllowedIPs = val

	return nil
}
