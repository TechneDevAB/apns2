// +build appengine

package apns2

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"google.golang.org/appengine/socket"
)

type GAETransport struct {
	httpClient  *http.Client
	certificate *tls.Certificate
	GConn       *socket.Conn
	Ctx         context.Context
}

func NewGAETransport(certificate tls.Certificate) *GAETransport {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	if len(certificate.Certificate) > 0 {
		tlsConfig.BuildNameToCertificate()
	}

	gaeTransport := &GAETransport{}

	transport := &http2.Transport{
		TLSClientConfig: tlsConfig,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			//gConn, err := socket.Dial(gaeTransport.Ctx, network, addr)
			//if err != nil {
			//	return nil, err
			//}

			timeout := time.Minute
			gConn, err := socket.DialTimeout(gae.ctx, "tcp", addr, timeout)
			if err != nil {
				return nil, err
			}

			gaeTransport.GConn = gConn

			tlsConn := tls.Client(gConn, cfg)
			return tlsConn, nil
		},
	}

	gaeTransport.httpClient = &http.Client{
		Transport: transport,
		Timeout:   HTTPClientTimeout,
	}

	return gaeTransport
}

// SetContext assigns a new context to the underlying socket.Conn
//  client := NewGAEClient(cert)
//  client.SetContext(ctx)
//  client.Push(notification)
func (c *GAETransport) SetContext(ctx context.Context) {
	c.Ctx = ctx
	if c.GConn != nil {
		c.GConn.SetContext(c.Ctx)
	}
}
func (t *GAETransport) HTTPClient() *http.Client {
	return t.httpClient
}

func (t *GAETransport) Certificate() *tls.Certificate {
	return t.certificate
}
