package Communication

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"goclient/K8s"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// server host url
	url string = "192.168.99.112"

	// kinds
	kb string = "kb"
	pg string = "pg"
	es string = "es"
)

type Comm struct {
	Client *http.Client
}

// general method to do request (curl -X <method> -u <username:password> -url <url:port/path> -d data)
func (comm *Comm) Curl(kind, username, password, path, method string, port int32, data []byte) (*http.Response, error) {
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("https://%s:%s@%s:%d/%s", username, password, url, port, path),
		bytes.NewBuffer(data),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if kind == "kb" {
		req.Header.Set("kbn-xsrf", "true")
	}
	resp, err := comm.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// create a http client, which is able to communicate with https server
func CreateHttpClient(k8sApi *K8s.K8sApi, kind, namespace, name string) (*http.Client, error) {
	cert, key, err := getCert(k8sApi, kind, namespace, name)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{pair},
				InsecureSkipVerify: true,
			},
		},
	}
	return client, nil
}

func getCert(k8sApi *K8s.K8sApi, kind, namespace, name string) ([]byte, []byte, error) {
	var secret *v1.Secret
	var err error
	if kind == "pg" {
		secret, err = k8sApi.GetSecret(metav1.GetOptions{}, namespace, "pgo.tls")
	} else {
		secret, err = k8sApi.GetSecret(metav1.GetOptions{}, namespace, name+"-"+kind+"-http-certs-internal")
	}
	if err != nil {
		logrus.Error(err)
		return nil, nil, err
	}
	return secret.Data["tls.crt"], secret.Data["tls.key"], nil
}

func stringifyResponse(resp *http.Response, err error) (string, error) {
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
