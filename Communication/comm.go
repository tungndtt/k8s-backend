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
)

var (
	// server host url
	url string = "ganmo.com"

	// kinds
	kb string = "kb"
	pg string = "pg"
	es string = "es"
)

type Comm struct {
	Client *http.Client
	Path   string
}

// general method to do request (curl -X <method> -u <username:password> -url <url:port/path> -d data)
func (comm *Comm) Curl(kind, username, password, path, method string, data []byte) (*http.Response, error) {
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("https://%s:%s@%s/%s/%s", username, password, url, comm.Path, path),
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
func getCommunication(kubeconfig, kind, namespace, name string) (*Comm, error) {
	api := Api{Kubeconfig: kubeconfig}
	secretName := ""
	if kind == kb {
		kbApi, err := api.KibanaAPI()
		if err != nil {
			return nil, err
		}
		kibana, err := kbApi.Get(namespace, name)
		if err != nil {
			return nil, err
		}
		secretName = kibana.Spec.Http.Tls.Cert.Secret
	} else if kind == es {
		esApi, err := api.ElasticsearchAPI()
		if err != nil {
			return nil, err
		}
		elastic, err := esApi.Get(namespace, name)
		if err != nil {
			return nil, err
		}
		secretName = elastic.Spec.Http.Tls.Cert.Secret
	}
	k8sApi, err := api.K8sAPI()
	if err != nil {
		return nil, err
	}
	cert, key, err := getCert(k8sApi, kind, namespace, name, secretName)
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

	var path string
	if kind == pg {
		path = namespace + "/postgres-operator"
	} else if kind == kb || kind == es {
		path = fmt.Sprintf("%s/%s-%s-http", namespace, name, kind)
	}
	comm := Comm{
		Client: client,
		Path:   path,
	}
	return &comm, nil
}

func getCert(k8sApi *K8s.K8sApi, kind, namespace, name, secretname string) ([]byte, []byte, error) {
	var secret *v1.Secret
	var err error
	if secretname != "" {
		secret, err = k8sApi.GetSecret(namespace, secretname)
	} else {
		if kind == pg {
			secret, err = k8sApi.GetSecret(namespace, "pgo.tls")
		} else if kind == kb || kind == es {
			secret, err = k8sApi.GetSecret(namespace, name+"-"+kind+"-http-certs-internal")
		}
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
