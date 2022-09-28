package main

import (
	"flag"
	"github.com/igor-kupczynski/fips-echo-server/echo"
	"github.com/igor-kupczynski/fips-echo-server/parsetls"
	"log"
)

var (
	address    = flag.String("address", "localhost:8443", "address for the server to listen on")
	certFile   = flag.String("certFile", "certs/domain.pem", "path to server certificate")
	keyFile    = flag.String("keyFile", "certs/domain.key", "path to server key")
	tlsVersion = flag.String("tlsVersion", "", "min TLS version, e.g. TLSv1.3")
	tlsCiphers = flag.String("tlsCiphers", "", "ciphersuites to use in mozilla string format, e.g. TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256:ECDHE-RSA-AES128-GCM-SHA256")
)

func main() {
	flag.Parse()

	parsedTlsVersion, parsedCiphers, err := parseTlsConfig()
	if err != nil {
		log.Fatal(err)
	}

	srv := echo.NewServer(*address, *certFile, *keyFile, parsedTlsVersion, parsedCiphers)

	ready := make(chan struct{})
	go func(ready <-chan struct{}) {
		<-ready
		log.Printf("Listening on https://%s with cert=%s and key=%s\n", *address, *certFile, *keyFile)
	}(ready)
	log.Fatal(srv.Serve(ready))
}

func parseTlsConfig() (version *uint16, ciphers []uint16, err error) {
	if *tlsVersion != "" {
		v, err := parsetls.Version(*tlsVersion)
		if err != nil {
			return nil, nil, err
		}
		version = &v
	}

	if *tlsCiphers != "" {
		ciphers, err = parsetls.CipherSuites(*tlsCiphers)
		if err != nil {
			return nil, nil, err
		}
	}

	return
}
