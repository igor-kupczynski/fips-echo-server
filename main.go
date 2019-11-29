package main

import (
	"flag"
	"github.com/igor-kupczynski/fips-echo-server/echo"
	"log"
)

var port = flag.Int("port", 8080, "port for the server to listen on")

func main() {
	flag.Parse()

	config := echo.Config{
		Port: *port,
	}

	ready := make(chan struct{})
	go func(ready <-chan struct{}) {
		<-ready
		log.Printf("Listening on %s\n", config.Address())
	}(ready)
	log.Fatal(echo.Serve(config, ready))
}
