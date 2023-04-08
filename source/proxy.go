package source

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Proxy struct{}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		p.handleHTTPS(w, r)
	} else {
		p.handleHTTP(w, r)
	}
}

func (p *Proxy) handleHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Proxying HTTPS request for %s://%s\n", r.Method, r.Host)
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   r.Host,
	})
	proxy.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	proxy.ServeHTTP(w, r)
}

func (p *Proxy) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	log.Printf("Proxying HTTPS request for %s://%s\n", r.Method, r.Host)
	hij, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	conn, _, err := hij.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer conn.Close()

	host := r.URL.Host
	if _, _, err := net.SplitHostPort(host); err != nil {
		host = net.JoinHostPort(host, "443")
	}

	remote, err := net.DialTimeout("tcp", host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer remote.Close()

	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\n\r\n")
	go io.Copy(remote, conn)
	io.Copy(conn, remote)
}

func ProxyMain() {
	proxy := &Proxy{}
	server := &http.Server{
		Addr:    ":18080",
		Handler: proxy,
	}
	fmt.Println("Proxy server listening on port 18080...")
	server.ListenAndServe()
}
