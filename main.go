package main

import (
	"flag"
	"github.com/igor-kupczynski/fips-echo-server/echo"
	"log"
)

var (
	address  = flag.String("address", "localhost:8443", "address for the server to listen on")
	certFile = flag.String("certFile", "certs/domain.pem", "path to server certificate")
	keyFile  = flag.String("keyFile", "certs/domain.key", "path to server key")
)

func main() {
	flag.Parse()

	srv := echo.Server(*address, *certFile, *keyFile)

	ready := make(chan struct{})
	go func(ready <-chan struct{}) {
		<-ready
		log.Printf("Listening on https://%s with cert=%s and key=%s\n", *address, *certFile, *keyFile)
	}(ready)
	log.Fatal(srv.Serve(ready))
}
