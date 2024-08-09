package utils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/atompi/activemq_exporter/options"
)

func FetchData(url string, data *[]byte) error {
	opts := options.Opts

	hc := &http.Client{}
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	r.Close = true

	if strings.HasPrefix(opts.Jolokia.URL, "https://") {
		clientTLSCert, err := tls.LoadX509KeyPair(opts.Jolokia.Cert, opts.Jolokia.Key)
		if err != nil {
			return err
		}
		RootCAPool, err := x509.SystemCertPool()
		if err != nil {
			return err
		}
		if caCert, err := os.ReadFile(opts.Jolokia.CA); err != nil {
			return err
		} else if ok := RootCAPool.AppendCertsFromPEM(caCert); !ok {
			return err
		}
		tlsConfig := &tls.Config{
			RootCAs: RootCAPool,
			Certificates: []tls.Certificate{
				clientTLSCert,
			},
			InsecureSkipVerify: opts.Jolokia.InsecureSkipVerify,
		}
		tr := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		hc = &http.Client{Transport: tr}
	}

	if opts.Jolokia.BasicAuth {
		auth := opts.Jolokia.Username + ":" + opts.Jolokia.Password
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		r.Header.Set("Authorization", basicAuth)
	}

	resp, err := hc.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("get jolokia api failed, status code: %d", resp.StatusCode)
		return err
	}

	buf := &bytes.Buffer{}
	buf.ReadFrom(resp.Body)
	*data = buf.Bytes()

	return nil
}
