package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var (
	httpAddr  = flag.String("http", ":80", "address to listen on")
	httpsAddr = flag.String("https", ":443", "address to listen on")
	source    = flag.String("source", "https://example.com", "source to proxy")
	crtStr    = flag.String("crt", "", "certificate in string")
	keyStr    = flag.String("key", "", "key in string")
)

func main() {
	flag.Parse()
	fmt.Printf("crt: %s", *crtStr)
	sourceUrl, err := url.Parse(*source)
	if err != nil {
		fmt.Printf("Error parsing source url %s: %s", *source, err)
		return
	}
	httpsHandle := http.NewServeMux()
	httpsHandle.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy := httputil.NewSingleHostReverseProxy(sourceUrl)
		proxy.ServeHTTP(w, r)
	}))

	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = make([]tls.Certificate, 1)
	tlsConfig.Certificates[0], err = tls.X509KeyPair([]byte(*crtStr), []byte(*keyStr))
	if err != nil {
		fmt.Printf("Error parsing certificate: %s", err)
		return
	}
	tlsServer := http.Server{Addr: *httpsAddr, TLSConfig: tlsConfig, Handler: httpsHandle}
	go func() {
		err = tlsServer.ListenAndServeTLS("", "")
		if err != nil {
			fmt.Printf("Error listening on %s: %s", *httpsAddr, err)
			os.Exit(1)
		}
	}()

	server := &http.Server{Addr: *httpAddr, Handler: httpsHandle}
	// http2.ConfigureServer(server, &http2.Server{IdleTimeout: time.Minute})
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %s", err)
		return
	}
}
