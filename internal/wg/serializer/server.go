package serializer

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"

	"github.com/fabiodcorreia/wg-concierge/internal/wg"
)

// Server Interface Properties
var regexServerAddresss = regexp.MustCompile(`Address\s*=\s*(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\/\d{2})`)
var regexServerPrivKey = regexp.MustCompile(`PrivateKey\s*=\s*(.{44})`)
var regexServerListenPort = regexp.MustCompile(`ListenPort\s*=\s*(\d+)`)

// Server Peers Properties
var regexServerPubKey = regexp.MustCompile(`PublicKey\s*=\s*(.{44})`)
var regexServerAllowIP = regexp.MustCompile(`AllowedIPs\s*=\s*(.+)`)

func UnmarshalServer(r io.Reader, sc *wg.ServerConfig) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("fail to unmarshal server configuration: %w", err)
	}

	return UnmarshalServerFromStr(string(buf), sc)
}

func UnmarshalServerFromStr(c string, cs *wg.ServerConfig) error {
	// Interface
	val, err := findConfig(&c, regexServerAddresss, true)
	if err != nil {
		return fmt.Errorf("server interface address: %w", err)
	}
	cs.Interface.Address = val

	val, err = findConfig(&c, regexServerPrivKey, true)
	if err != nil {
		return fmt.Errorf("server interface private key: %w", err)
	}
	cs.Interface.PrivateKey = val

	val, err = findConfig(&c, regexServerListenPort, true)
	if err != nil {
		return fmt.Errorf("server interface listen port : %w", err)
	}
	cs.Interface.PrivateKey = val

	// Peer
	// TODO

	return nil
}
