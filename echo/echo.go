// echo contains an echo server implementation
package echo

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

func init() {
	http.HandleFunc("/", echoHandler)
}

// Config is the echo server configuration
type Config struct {
	Port int
}

// Address returns address on which the server listens
func (c *Config) Address() string {
	return fmt.Sprintf(":%d", c.Port)
}

const echoLimit = 140

func echoHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, echoLimit)
	n, _ := r.Body.Read(b)
	defer r.Body.Close()
	if n > 0 {
		_, _ = w.Write(b[:n])
	}
}

var srv *http.Server = &http.Server{}

// Serve starts the echo server
func Serve(c Config, ready chan<- struct{}) error {

	ln, err := net.Listen("tcp", c.Address())
	if err != nil {
		return err
	}
	ready <- struct{}{}
	return srv.Serve(ln)
}

// Shutdown the running server
func Shutdown() error {
	return srv.Shutdown(context.Background())
}
