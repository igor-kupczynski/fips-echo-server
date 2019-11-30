package echo

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

const text140chars = "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"

var (
	caFile   = "../certs/ca.pem"
	certFile = "../certs/domain.pem"
	keyFile  = "../certs/domain.key"
)

// TestServe runs end-to-end smoke test of the echo server
func TestServe(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "Echo the message back to the client",
			in:   "hello",
			out:  "hello",
		},
		{
			name: "Limit to 140 characters",
			in:   fmt.Sprintf("%s%s", text140chars, "and even more"),
			out:  text140chars,
		},
	}
	port := 16123
	for _, tt := range tests {
		port++
		t.Run(tt.name, func(t *testing.T) {
			s := Server(fmt.Sprintf("localhost:%d", port), certFile, keyFile, nil, nil)
			ready := make(chan struct{})
			go func(ready <-chan struct{}, in, out string) {
				<-ready

				client := buildHttpsClient(t, caFile)
				target := fmt.Sprintf("https://%s", s.address)
				resp, err := client.Post(target, "text/plain", strings.NewReader(in))
				if err != nil {
					t.Errorf("Can't connect to server = %v", err)
					return
				}
				defer resp.Body.Close()

				bytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("Can't read the response = %v", err)
				} else if got := string(bytes); got != out {
					t.Errorf("got %s, but want %s", got, out)
				}

				if err := s.Shutdown(); err != nil {
					t.Errorf("Error shutting down = %v", err)
					return
				}
			}(ready, tt.in, tt.out)
			if err := s.Serve(ready); err != nil && !strings.Contains(err.Error(), "http: Server closed") {
				t.Errorf("Serve() error = %v", err)
			}
		})
	}
}

func buildHttpsClient(t *testing.T, caFile string) *http.Client {
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		t.Errorf("Can't load CA cert err=%v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{Transport: transport}
}
