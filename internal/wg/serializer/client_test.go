package serializer_test

import (
	"testing"

	"github.com/fabiodcorreia/wg-concierge/internal/wg"
	"github.com/fabiodcorreia/wg-concierge/internal/wg/serializer"
)

func TestMarshalClientToStr(t *testing.T) {
	type args struct {
		cc wg.ClientConfig
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
		wantErr    bool
	}{
		{
			name: "",
			args: args{wg.ClientConfig{Interface: wg.ClientInterface{
				Address:    "",
				DNS:        "",
				PrivateKey: "",
			}, Peer: wg.ClientPeer{
				AllowedIPs: "",
				Endpoint:   "",
				KeepAlive:  11,
				PublicKey:  "",
			}}},
			wantResult: `
			[Interface]
			[Peer]
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := serializer.MarshalClientToStr(tt.args.cc)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalClientToStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("MarshalClientToStr() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
