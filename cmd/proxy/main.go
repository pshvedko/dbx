package main

import (
	"bufio"
	"crypto/tls"
	_ "embed"
	"fmt"
	"log"
	"net"
	"net/http"
)

//go:embed pem/cert.pem
var cert []byte

//go:embed pem/key.pem
var key []byte

type Proxy struct {
	tls.Certificate
}

func (p Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", r)
	switch r.Method {
	case http.MethodConnect:
		conn3, err := tls.DialWithDialer(&net.Dialer{Cancel: r.Context().Done()}, "tcp", r.Host, &tls.Config{})
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			log.Print(err)
			return
		}
		defer conn3.Close()

		conn1, _, err := http.NewResponseController(w).Hijack()
		if err != nil {
			w.WriteHeader(http.StatusTeapot)
			log.Print(err)
			return
		}

		fmt.Fprintln(conn1, "HTTP/1.1 200 Connection established")
		fmt.Fprintln(conn1)

		conn2 := tls.Server(conn1, &tls.Config{
			GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
				return &p.Certificate, nil
			},
		})
		defer conn2.Close()

		err = conn2.HandshakeContext(r.Context())
		if err != nil {
			log.Print(err)
			return
		}

		r2 := bufio.NewReader(conn2)
		r3 := bufio.NewReader(conn3)
		for {
			req, err := http.ReadRequest(r2)
			if err != nil {
				log.Print(err)
				return
			}

			req.Write(conn3)

			res, err := http.ReadResponse(r3, req)
			if err != nil {
				log.Print(err)
				return
			}

			res.Write(conn2)
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
	err = http.ListenAndServe(":8080", Proxy{Certificate: certificate})
	if err != nil {
		log.Print(err)
		return
	}
}
