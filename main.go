package main

import (
	"flag"
	"github.com/igor-kupczynski/fips-echo-server/echo"
	"log"
)

var (
	port     = flag.Int("port", 8443, "port for the server to listen on")
	certFile = flag.String("certFile", "certs/domain.pem", "path to server certificate")
	keyFile  = flag.String("keyFile", "certs/domain.key", "path to server key")
)

func main() {
	flag.Parse()

	srv := echo.Server(*port, *certFile, *keyFile)

	ready := make(chan struct{})
	go func(ready <-chan struct{}) {
		<-ready
		log.Printf("Listening on https://localhost:%d with cert=%s and key=%s\n", *port, *certFile, *keyFile)
	}(ready)
	log.Fatal(srv.Serve(ready))
}
