// Package parsetls contains utils for parsing the TLS settings.
//
// We need this because golang stdlib contains only numeric constants for the
// ciphers and the protocols.
package parsetls

import (
	"crypto/tls"
	"fmt"
	"strings"
)

var protocols = map[string]uint16{
	"TLSv1.0": tls.VersionTLS10,
	"TLSv1.1": tls.VersionTLS11,
	"TLSv1.2": tls.VersionTLS12,
	"TLSv1.3": tls.VersionTLS13,
}

// Version returns the numeric constant representing the TLS version
// corresponding to the provided string, e.g. "TLSv1.3" --> tls.VersionTLS13.
func Version(input string) (uint16, error) {
	if v, ok := protocols[input]; ok {
		return v, nil
	}
	return 0, fmt.Errorf("non existing protocol %s", input)
}

var ciphers = map[string]uint16{
	"AES128-GCM-SHA256":             tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	"AES128-SHA":                    tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	"AES128-SHA256":                 tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
	"AES256-GCM-SHA384":             tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	"AES256-SHA":                    tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	"DES-CBC3-SHA":                  tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
	"ECDHE-ECDSA-AES128-GCM-SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	"ECDHE-ECDSA-AES128-SHA":        tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	"ECDHE-ECDSA-AES128-SHA256":     tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
	"ECDHE-ECDSA-AES256-GCM-SHA384": tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	"ECDHE-ECDSA-AES256-SHA":        tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	"ECDHE-ECDSA-CHACHA20-POLY1305": tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	"ECDHE-RSA-AES128-GCM-SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"ECDHE-RSA-AES128-SHA":          tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	"ECDHE-RSA-AES128-SHA256":       tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	"ECDHE-RSA-AES256-GCM-SHA384":   tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	"ECDHE-RSA-AES256-SHA":          tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	"ECDHE-RSA-CHACHA20-POLY1305":   tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	"ECDHE-RSA-DES-CBC3-SHA":        tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	"TLS_AES_128_GCM_SHA256":        tls.TLS_AES_128_GCM_SHA256,
	"TLS_AES_256_GCM_SHA384":        tls.TLS_AES_256_GCM_SHA384,
	"TLS_CHACHA20_POLY1305_SHA256":  tls.TLS_CHACHA20_POLY1305_SHA256,
}

// CipherSuites returns a slice of numeric constants corresponding to the list
// provided as a string. The list should be ":" delimited.
// E.g. "TLS_CHACHA20_POLY1305_SHA256:TLS_AES_256_GCM_SHA384" -->
// { tls.TLS_CHACHA20_POLY1305_SHA256, tls.TLS_AES_256_GCM_SHA384 }.
func CipherSuites(input string) ([]uint16, error) {
	items := strings.Split(input, ":")
	cipers := make([]uint16, len(items))
	for i := 0; i < len(items); i++ {
		c, ok := ciphers[items[i]]
		if !ok {
			return nil, fmt.Errorf("non existing cipher %s", items[i])
		}
		cipers[i] = c
	}
	return cipers, nil
}
