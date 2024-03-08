package request

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/dop251/goja"
)

var (
	// NOTE: before calling any request functions, tlsConfig must be initialized.
	tlsConfig *tls.Config
)

func InitTrustedHttpsCerts(certs []string) error {
	var err error
	tlsConfig, err = LoadTlsConfig(certs)
	return err
}

type HttpResponse struct {
	Status     string
	StatusCode int
	Headers    [][2]string
	Body       string
}

func HttpsRequest(method, serverURL, body string, headers ...string) HttpResponse {
	if serverURL == "" {
		panic(goja.NewSymbol("Empty url"))
	}

	req, err := newHttpRequest(method, serverURL, body, headers...)
	if err != nil {
		panic(goja.NewSymbol("Error in parsing http request: " + err.Error()))
	}
	// disable keepalive
	req.Close = true

	client := &http.Client{
		Transport: &http.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				conn, err := tls.Dial(network, addr, tlsConfig)
				return conn, err
			},
		},
	}

	resp, err := client.Do(&req)
	if err != nil {
		panic(goja.NewSymbol("Error in sending http request: " + err.Error()))
	}
	result, err := newHttpResponse(resp)
	if err != nil {
		panic(goja.NewSymbol("Error in parsing http response: " + err.Error()))
	}
	return result
}

func newHttpRequest(method, serverURL, body string, headers ...string) (result http.Request, err error) {
	result.Method = method
	result.URL, err = url.Parse(serverURL)
	if err != nil {
		return
	}
	result.Header = make(http.Header)
	for _, h := range headers {
		fields := strings.Split(h, ":")
		if len(fields) != 2 {
			return result, errors.New("Invalid header: " + h)
		}
		result.Header.Add(strings.TrimSpace(fields[0]), strings.TrimSpace(fields[1]))
	}
	result.Body = io.NopCloser(strings.NewReader(body))
	return
}

func newHttpResponse(resp *http.Response) (result HttpResponse, err error) {
	result.Status = resp.Status
	result.StatusCode = resp.StatusCode
	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return
	}
	result.Body = buf.String()
	keys := make([]string, 0, len(resp.Header))
	for k := range resp.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	result.Headers = make([][2]string, 0, len(keys))
	for _, k := range keys {
		for _, v := range resp.Header[k] {
			result.Headers = append(result.Headers, [2]string{k, v})
		}
	}
	return
}

func LoadTlsConfig(certs []string) (*tls.Config, error) {
	tlsConfig := &tls.Config{RootCAs: x509.NewCertPool()}

	for _, c := range certs {
		var b *pem.Block
		b, _ = pem.Decode([]byte(c))
		x509Cert, err := x509.ParseCertificate(b.Bytes)
		if err != nil {
			return nil, err
		}
		tlsConfig.RootCAs.AddCert(x509Cert)
	}
	return tlsConfig, nil
}
