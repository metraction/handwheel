package integrations

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"

	"github.com/metraction/handwheel/model"
)

func NewHttpTransport(cfg *model.Config) *http.Transport {
	var caCertPool *x509.CertPool
	var err error
	caCertPool, err = x509.SystemCertPool()
	if err != nil {
		log.Printf("failed to load system cert pool: %v", err)
	}

	if cfg.CAFile != "" {
		func() {
			log.Println("Using CA root file.", cfg.CAFile)
			pem, err := os.ReadFile(cfg.CAFile)
			if err != nil {
				log.Printf("failed to read CA root: %v", err)
				return
			}
			if !caCertPool.AppendCertsFromPEM(pem) {
				log.Printf("failed to append CA root cert")
				return
			}
		}()
	} else if cfg.CARootPEM != "" {
		log.Println("Using CA root PEM.")
		pemBytes := normalizePEM([]byte(cfg.CARootPEM))
		if !caCertPool.AppendCertsFromPEM(pemBytes) {
			log.Printf("failed to append CA root cert")
		}
	}
	return &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: caCertPool},
	}
}
