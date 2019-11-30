package parsetls

import (
	"crypto/tls"
	"reflect"
	"testing"
)

func TestCipherSuites(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    []uint16
		wantErr bool
	}{
		{
			"Modern",
			"TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256:ECDHE-RSA-AES128-GCM-SHA256",
			[]uint16{tls.TLS_AES_128_GCM_SHA256, tls.TLS_AES_256_GCM_SHA384, tls.TLS_CHACHA20_POLY1305_SHA256, tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
			false,
		},
		{
			"Intermediate v1.2 ciphers",
			"ECDHE-ECDSA-AES128-GCM-SHA256" + ":" + "ECDHE-RSA-AES128-GCM-SHA256" + ":" +
				"ECDHE-ECDSA-AES256-GCM-SHA384" + ":" + "ECDHE-RSA-AES256-GCM-SHA384" + ":" +
				"ECDHE-ECDSA-CHACHA20-POLY1305" + ":" + "ECDHE-RSA-CHACHA20-POLY1305",
			[]uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			},
			false,
		},
		{
			"Invalid",
			"ECDHE-ECDSA-AES128-GCM-SHA256:foobar",
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CipherSuites(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("CipherSuites() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CipherSuites() got = %v, want %v", got, tt.want)
			}
		})

	}
}

func TestVersion(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    uint16
		wantErr bool
	}{
		{"TLS v1.0", "TLSv1.0", tls.VersionTLS10, false},
		{"TLS v1.1", "TLSv1.1", tls.VersionTLS11, false},
		{"TLS v1.2", "TLSv1.2", tls.VersionTLS12, false},
		{"TLS v1.3", "TLSv1.3", tls.VersionTLS13, false},
		{"Invalid", "TLSv1.9", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Version(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Version() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Version() got = %v, want %v", got, tt.want)
			}
		})
	}
}
