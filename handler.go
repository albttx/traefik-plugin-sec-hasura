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
	GraphQLPath   string   `json:"graphql_path,omitempty"`
	IgnoreHeaders []string `json:"ignore_headers"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		GraphQLPath:   "/v1/graphql",
		IgnoreHeaders: []string{},
	}
}

// HasuraPlugin a HasuraPlugin plugin.
type HasuraPlugin struct {
	cfg Config

	next    http.Handler
	headers map[string]string
	name    string
}

// New created a new HasuraPlugin plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &HasuraPlugin{
		cfg: *config,

		next: next,
		name: name,
	}, nil
}

func (p *HasuraPlugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	os.Stderr.WriteString("ServeHTTP ========================== " + req.URL.Path + "\n")

	// ignore if it's not graphql endpoint
	if req.URL.Path != p.cfg.GraphQLPath {
		p.next.ServeHTTP(rw, req)
		return
	}

	// ignore if some headers are specify
	for _, h := range p.cfg.IgnoreHeaders {
		if req.Header.Get(h) != "" {
			p.next.ServeHTTP(rw, req)
			return
		}
	}

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

	p.next.ServeHTTP(rw, req)
}
