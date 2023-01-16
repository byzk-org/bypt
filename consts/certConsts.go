package consts

import (
	_ "embed"
)

var (
	//go:embed certs/ca.cert
	CaCert []byte
	//go:embed certs/user.cert
	UserCert []byte
	//go:embed certs/user.key
	UserKey []byte
)

