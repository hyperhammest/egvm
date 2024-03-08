package request

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"
)

const (
	TrustedCertsPathForTest = "./certs"
)

const (
	HttpScriptTemplate = `
		//const resp = HttpsRequest('GET', 'https://elfinauth.paralinker.io/smartbch/eh_ping', '', 'Content-Type:application/json')
		const resp = HttpsRequest('GET', 'https://elfincdn111.paralinker.io/eh_ping', '', 'Content-Type:application/json')
		const body = resp.Body
	`
)

func setupGojaVmForHttp() *goja.Runtime {
	vm := goja.New()
	vm.Set("HttpsRequest", HttpsRequest)
	return vm
}

func loadTlsConfigForTest(certsPath string) (*tls.Config, error) {
	dicEntry, err := os.ReadDir(certsPath)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{RootCAs: x509.NewCertPool()}
	for _, e := range dicEntry {
		if e.IsDir() {
			panic("cannot keep dir in trusted certs directory")
		}

		fileNames := strings.Split(e.Name(), ".")

		// if not pem file
		if len(fileNames) != 2 || fileNames[1] != "pem" {
			continue
		}

		certBz, err := os.ReadFile(filepath.Join(certsPath, e.Name()))
		if err != nil {
			return nil, err
		}

		var b *pem.Block
		for len(certBz) > 0 {
			b, certBz = pem.Decode(certBz)
			x509Cert, err := x509.ParseCertificate(b.Bytes)
			if err != nil {
				return nil, err
			}
			tlsConfig.RootCAs.AddCert(x509Cert)
		}

	}
	return tlsConfig, nil
}

func TestHttpRequest(t *testing.T) {
	vm := setupGojaVmForHttp()
	tlsConfig, _ = loadTlsConfigForTest(TrustedCertsPathForTest)
	_, err := vm.RunString(HttpScriptTemplate)
	require.NoError(t, err)

	resp := vm.Get("resp").Export().(HttpResponse)
	require.EqualValues(t, 200, resp.StatusCode)
	require.EqualValues(t, `{"isSuccess":true,"message":"pong"}`, resp.Body)
}
