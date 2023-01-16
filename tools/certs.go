package tools

import (
	"encoding/pem"
	"github.com/tjfoc/gmsm/x509"
	"time"
)

func ValidCertExpire(certPemBytes []byte) bool {
	certificate, err := ParseCertByPem(certPemBytes)
	if err != nil {
		return false
	}
	return time.Now().Before(certificate.NotAfter)
}

func ParseCertByPem(certPemBytes []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certPemBytes)
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return certificate, nil
}


