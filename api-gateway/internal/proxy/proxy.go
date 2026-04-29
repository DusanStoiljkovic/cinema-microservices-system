package proxy

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ProxyHandler struct {
	routes map[string]*url.URL
	client *http.Client
}

func NewProxyHandler(timeout time.Duration) *ProxyHandler {
	return &ProxyHandler{
		routes: make(map[string]*url.URL),
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (p *ProxyHandler) AddRoute(prefix string, backend string) error {
	u, err := url.Parse(backend)
	if err != nil {
		return err
	}
	p.routes[prefix] = u
	return nil
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var target *url.URL
	var matchedPrefix string

	// FIND ROUTE
	for prefix, backend := range p.routes {
		if strings.HasPrefix(r.URL.Path, prefix) && len(prefix) > len(matchedPrefix) {
			target = backend
			matchedPrefix = prefix
		}
	}

	if target == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	// CREATE NEW REQUEST
	outReq := p.createRequest(r, target, matchedPrefix)
	if outReq == nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	// SEND REQUEST
	resp, err := p.client.Do(outReq)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		log.Println("Proxy Error: ", err)
		return
	}
	defer resp.Body.Close()

	// COPY HEADERS
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	// COPY BODY
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println("copy error: ", err)
	}

	log.Printf("%s %s -> %s (%v)",
		r.Method,
		r.URL.Path,
		target.Host,
		time.Since(start),
	)
}

func (p *ProxyHandler) createRequest(r *http.Request, target *url.URL, prefix string) *http.Request {
	outURL := *r.URL
	outURL.Scheme = target.Scheme
	outURL.Host = target.Host

	path := strings.TrimPrefix(r.URL.Path, "/api")

	outURL.Path = path

	req, err := http.NewRequestWithContext(
		r.Context(),
		r.Method,
		outURL.String(),
		r.Body,
	)
	if err != nil {
		log.Print("Request creation error: ", err)
		return nil
	}

	// COPY HEADERS
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// FORWARDED HEADERS
	req.Header.Set("X-Forwarded-For", r.RemoteAddr)
	req.Header.Set("X-Forwarded-Host", r.Host)

	return req
}
