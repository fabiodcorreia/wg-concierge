package wg

// ClientInterface is the WireGuard Client Interface specification
type ClientInterface struct {
	PrivateKey string // Client Private Key
	Address    string // Client Virtual Private IP
	DNS        string // DNS Server
}

// ClientPeer is the WireGuard Client Peer specification
type ClientPeer struct {
	PublicKey  string // Server Public Key
	Endpoint   string // Server Endpoint
	AllowedIPs string // IP ranges routed to the VPN
	KeepAlive  uint32 // Seconds between each keep alive
}

// ClientConfig is the WireGuard Client Configuration specification
type ClientConfig struct {
	Interface ClientInterface // Client Interface specification
	Peer      ClientPeer      // Client Peer specificiation
}

// ServerInterface is the WireGuard Server Interface specification
type ServerInterface struct {
	PrivateKey string
	PublicKey  string
	Address    string
	SaveConfig bool
	ListenPort uint32
}

// ServerPeer is the WireGuard Server Peer specification
type ServerPeer struct {
	PublicKey  string
	AllowedIPs string
}

// ServerConfig is the WireGuard Server Configuration specification
type ServerConfig struct {
	Interface ServerInterface
	Peers     []ServerPeer
}
