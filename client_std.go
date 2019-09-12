// +build local

package apns2

import (
	"crypto/tls"
	"golang.org/x/net/http2"
	"net"
	"net/http"
)

type StdTransport struct {
	httpClient  *http.Client
	certificate *tls.Certificate
}

// NewClient returns a new Client with an underlying http.Client configured with
// the correct APNs HTTP/2 transport settings. It does not connect to the APNs
// until the first Notification is sent via the Push method.
//
// As per the Apple APNs Provider API, you should keep a handle on this client
// so that you can keep your connections with APNs open across multiple
// notifications; don’t repeatedly open and close connections. APNs treats rapid
// connection and disconnection as a denial-of-service attack.
//
// If your use case involves multiple long-lived connections, consider using
// the ClientManager, which manages clients for you.
func NewStdTransport(certificate tls.Certificate) *StdTransport {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	if len(certificate.Certificate) > 0 {
		tlsConfig.BuildNameToCertificate()
	}

	transport := &http2.Transport{
		TLSClientConfig: tlsConfig,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return tls.DialWithDialer(&net.Dialer{Timeout: TLSDialTimeout}, network, addr, cfg)
		},
	}

	return &StdTransport{
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   HTTPClientTimeout,
		},
	}
}

func (t *StdTransport) HTTPClient() *http.Client {
	return t.httpClient
}

func (t *StdTransport) Certificate() *tls.Certificate {
	return t.certificate
}
