package netutils

import (
	"crypto/tls"
	"fmt"
)

func LoadTlsCreds(cert, key string) (*tls.Config, error) {
	creds, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to load server cert and key: %w", err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{creds},
	}, nil
}
