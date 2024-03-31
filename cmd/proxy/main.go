package main

import (
	"bufio"
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

//go:embed pem/cert.pem
var cert []byte

//go:embed pem/key.pem
var key []byte

type Proxy struct {
	tls.Certificate
	tls.Dialer
	http.Client
}

func (p Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== %+v", r)

	conn1, _, err := http.NewResponseController(w).Hijack()
	if err != nil {
		w.WriteHeader(http.StatusTeapot)
		log.Print(err)
		return
	}

	switch r.Method {
	case http.MethodConnect:
		conn3, err := p.DialContext(r.Context(), "tcp", r.Host)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			log.Print(err)
			return
		}
		defer conn3.Close()

		fmt.Fprintln(conn1, "HTTP/1.1 200 Connection established")
		fmt.Fprintln(conn1, "Proxy-Connection: close")
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

			log.Printf("--> %+v", req)

			req.Write(conn3)

			res, err := http.ReadResponse(r3, req)
			if err != nil {
				log.Print(err)
				return
			}

			log.Printf("<-- %+v", res)

			res.Write(conn2)
		}
	default:
		r1 := bufio.NewReader(conn1)
		for {
			res, err := p.Do(&http.Request{
				Method:           r.Method,
				URL:              r.URL,
				Proto:            r.Proto,
				ProtoMajor:       r.ProtoMajor,
				ProtoMinor:       r.ProtoMinor,
				Header:           r.Header,
				Body:             r.Body,
				GetBody:          r.GetBody,
				ContentLength:    r.ContentLength,
				TransferEncoding: r.TransferEncoding,
				Close:            r.Close,
				Host:             r.Host,
				Trailer:          r.Trailer,
			})
			if err != nil {
				log.Print(err)
				return
			}

			log.Printf("<-- %+v", res)

			res.Write(conn1)

			r, err = http.ReadRequest(r1)
			if err != nil {
				log.Print(err)
				return
			}

			log.Printf("--> %+v", r)
		}
	}
}

func main() {
	certificate, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Print(err)
		return
	}

	s := http.Server{
		Addr:    ":8080",
		Handler: Proxy{Certificate: certificate},
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	context.AfterFunc(ctx, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err := s.Shutdown(ctx)
		if err != nil {
			log.Print(err)
			return
		}
	})

	err = s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Print(err)
		return
	}
}
