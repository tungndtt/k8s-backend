package Communication

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
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
	Client   *http.Client
	Path     string
	User     string
	Password string
}

// general method to do request (curl -X <method> -u <username:password> -url <url:port/path> -d data)
func (comm *Comm) Curl(kind, path, method string, data []byte) (*http.Response, error) {
	tmp := fmt.Sprintf("https://%s/%s/%s", url, comm.Path, path)
	fmt.Println(tmp)
	req, err := http.NewRequest(
		method,
		tmp,
		bytes.NewBuffer(data),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", comm.User, comm.Password))))
	if kind == "kb" {
		req.Header.Set("kbn-xsrf", "true")
	}
	resp, err := comm.Client.Do(req)
	if err != nil {
		return nil, err
	}
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(resp)
	fmt.Println(string(b))
	return resp, nil
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
