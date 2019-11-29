package echo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

const text140chars = "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"

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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{Port: 16123}
			ready := make(chan struct{})
			go func(ready <-chan struct{}, in, out string) {
				<-ready

				target := fmt.Sprintf("http://localhost:%d", c.Port)
				resp, err := http.Post(target, "text/plain", strings.NewReader(in))
				if err != nil {
					t.Errorf("Can't connect to server = %v", err)
					return
				}
				defer resp.Body.Close()

				bytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("Can't read the response = %v", err)
					return
				}

				if got := string(bytes); got != out {
					t.Errorf("got %s, but want %s", got, out)
				}

				if err := Shutdown(); err != nil {
					t.Errorf("Error shutting down = %v", err)
					return
				}
			}(ready, tt.in, tt.out)
			if err := Serve(c, ready); err != nil && !strings.Contains(err.Error(), "http: Server closed") {
				t.Errorf("Serve() error = %v", err)
			}
		})
	}
}
