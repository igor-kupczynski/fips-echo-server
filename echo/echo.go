// echo contains an echo server implementation
package echo

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

const echoLimit = 140

func echoHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, echoLimit)
	n, _ := r.Body.Read(b)
	defer r.Body.Close()
	if n > 0 {
		_, _ = w.Write(b[:n])
	}
}

// Server returns a server struct
func Server(port int, certFile, keyFile string) *server {
	mux := http.NewServeMux()
	return &server{
		port:     port,
		certFile: certFile,
		keyFile:  keyFile,
		mux:      mux,
		srv: &http.Server{
			Handler: mux,
		},
	}
}

// server is the echo server internal state
type server struct {
	// Port on which the server listens
	port int
	// Certificate file to use
	certFile string
	// Private key to the certificate
	keyFile string
	// Multiplexer to register the handles on
	mux *http.ServeMux
	// http.Server handle
	srv *http.Server
}

// address returns address on which the server listens
func (s *server) address() string {
	return fmt.Sprintf("localhost:%d", s.port)
}

// Serve starts the echo server
func (s *server) Serve(ready chan<- struct{}) error {

	s.mux.HandleFunc("/", echoHandler)

	ln, err := net.Listen("tcp", s.address())
	if err != nil {
		return err
	}
	ready <- struct{}{}
	return s.srv.ServeTLS(ln, s.certFile, s.keyFile)
}

// Shutdown the running server
func (s *server) Shutdown() error {
	return s.srv.Shutdown(context.Background())
}
