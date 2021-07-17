package qr_test

import (
	"bytes"
	"testing"

	"github.com/fabiodcorreia/wg-concierge/internal/qr"
)

func TestEncodeString(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "empty content",
			args: args{
				content: "",
			},
			want:    968,
			wantErr: false,
		},
		{
			name: "with content",
			args: args{
				content: "my content to QR Code",
			},
			want:    1048,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := qr.EncodeString(tt.args.content, w); (err != nil) != tt.wantErr {
				t.Errorf("EncodeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); len(gotW) != tt.want {
				t.Errorf("EncodeString() = %v, want %v", len(gotW), tt.want)
			}
		})
	}
}
