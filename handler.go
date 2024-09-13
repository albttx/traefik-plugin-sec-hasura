package traefik_plugin_sec_hasura

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// Config the plugin configuration.
type Config struct {
	Headers map[string]string `json:"headers,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers: make(map[string]string),
	}
}

// HasuraPLugin a HasuraPLugin plugin.
type HasuraPLugin struct {
	next    http.Handler
	headers map[string]string
	name    string
}

// New created a new HasuraPLugin plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &HasuraPLugin{
		headers: config.Headers,
		next:    next,
		name:    name,
	}, nil
}

func (a *HasuraPLugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	os.Stderr.WriteString("ServeHTTP ==========================\n")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		os.Stderr.WriteString(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	switch data.(type) {
	case []interface{}:
		http.Error(rw, "Batch queries are forbidden", http.StatusForbidden)
		return
	}

	a.next.ServeHTTP(rw, req)
}
