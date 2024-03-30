package main

import (
	"bufio"
	"crypto/tls"
	_ "embed"
	"fmt"
	"log"
	"net/http"
)

//go:embed pem/cert.pem
var cert []byte

//go:embed pem/key.pem
var key []byte

type Proxy struct {
	Certificates []tls.Certificate
}

func (p Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", r)
	switch r.Method {
	case http.MethodConnect:
		switch x := w.(type) {
		case http.Hijacker:
			conn, rw, err := x.Hijack()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Print(err)
				return
			}

			fmt.Fprintln(rw, "HTTP/1.1 200 Connection established")
			fmt.Fprintln(rw)
			rw.Flush()

			conn2 := tls.Server(conn, &tls.Config{
				Certificates: p.Certificates,
			})
			defer conn2.Close()

			err = conn2.HandshakeContext(r.Context())
			if err != nil {
				log.Print(err)
				return
			}

			r, _ := bufio.NewReader(conn2), bufio.NewWriter(conn2)
			for {
				req, err := http.ReadRequest(r)
				if err != nil {
					log.Print(err)
					return
				}
				log.Printf("%+v", req)

				fmt.Fprintln(conn2, "HTTP/1.1 200 OK")
				fmt.Fprintln(conn2, "Content-Length: 3")
				fmt.Fprintln(conn2)
				fmt.Fprintln(conn2, "OK")
			}
		default:
			w.WriteHeader(http.StatusTeapot)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	certificate, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Print(err)
		return
	}
	err = http.ListenAndServe(":8080", Proxy{Certificates: []tls.Certificate{certificate}})
	if err != nil {
		log.Print(err)
		return
	}
}
