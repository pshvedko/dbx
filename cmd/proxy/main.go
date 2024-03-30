package main

import (
	"bufio"
	"crypto/tls"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
		switch x := w.(type) {
		case http.Hijacker:
			conn, _, err := x.Hijack()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Print(err)
				return
			}

			fmt.Fprintln(conn, "HTTP/1.1 200 Connection established")
			fmt.Fprintln(conn)

			conn2 := tls.Server(conn, &tls.Config{
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

			r1 := bufio.NewReader(conn2)
			for {
				req, err := http.ReadRequest(r1)
				if err != nil {
					log.Print(err)
					return
				}
				log.Printf("%+v", req)

				res := http.Response{
					Status:        "OK",
					StatusCode:    http.StatusOK,
					Proto:         r.Proto,
					ProtoMajor:    r.ProtoMajor,
					ProtoMinor:    r.ProtoMinor,
					Body:          io.NopCloser(strings.NewReader("OK\n")),
					ContentLength: 3,
					Close:         true,
				}
				res.Write(conn2)
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
	err = http.ListenAndServe(":8080", Proxy{Certificate: certificate})
	if err != nil {
		log.Print(err)
		return
	}
}
