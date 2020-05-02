package main

import (
	"fmt"
	"os"

	"github.com/fabiodcorreia/wg-concierge/internal/wg"
	"github.com/fabiodcorreia/wg-concierge/internal/wg/serializer"
)

var version = "development"

/*
wg-concierge --wg-config=/etc/wireguard/wg0.conf --wg-store=/home/app/concierge.bin
*/

func main() {

	/*
		cc := wg.ClientConfig{
			Interface: wg.ClientInterface{
				Address:    "10.253.3.2/32",
				DNS:        "1.1.1.1",
				PrivateKey: "CHFXbC9IBiq+sUxilT5TSNFe9KnefDZmooOyMCGZA1U=",
			},
			Peer: wg.ClientPeer{
				PublicKey:  "badiMN09ahTH/kp6knl9ew708Leq39ii7NVeD+eedR4=",
				AllowedIPs: "0.0.0.0/0, ::/0",
				Endpoint:   "192.168.5.81:55555",
				KeepAlive:  60,
			},
		}
		s, _ := wg.MarshalClientToString(cc)

		fmt.Println(s)

		f, err := os.Create("./testdata/code.png")
		fmt.Println(err)
		defer f.Close()
		err = qr.EncodeString(s, f)
		fmt.Println(err)
	*/

	var cc wg.ClientConfig
	f, _ := os.OpenFile("./testdata/client.conf", os.O_RDONLY, os.ModePerm)
	defer f.Close()
	e := serializer.UnmarshalClient(f, &cc)
	fmt.Println(cc, e)
}
