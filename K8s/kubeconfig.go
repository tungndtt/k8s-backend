package K8s

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type KubeConfig struct {
	ApiVersion     string    `json:"apiVersion"`
	Kind           string    `json:"kind"`
	Clusters       []Cluster `json:"clusters"`
	Contexts       []Context `json:"contexts"`
	CurrentContext string    `json:"current-context"`
	Users          []User    `json:"users"`
}

type Cluster struct {
	ClusterInfo ClusterInfo `json:"cluster"`
	Name        string      `json:"name"`
}

type ClusterInfo struct {
	Cert      string `json:"certificate-authority"`
	ServerUrl string `json:"server"`
}

type Context struct {
	ContextInfo ContextInfo `json:"context"`
	Name        string      `json:"name"`
}

type ContextInfo struct {
	Cluster   string  `json:"cluster"`
	Namespace *string `json:"namespace,omitempty"`
	User      string  `json:"user"`
}

type User struct {
	Name     string   `json:"name"`
	UserInfo UserInfo `json:"user"`
}

type UserInfo struct {
	Token string `json:"token"`
}

func GenerateKubeconfig(token, user string) error {
	config := KubeConfig{
		ApiVersion: "v1",
		Kind:       "Config",
		Clusters: []Cluster{
			{
				Name: "minikube",
				ClusterInfo: ClusterInfo{
					Cert:      "/home/tung/.minikube/ca.crt",
					ServerUrl: "https://192.168.99.111:8443",
				},
			},
		},
		Contexts: []Context{
			{
				Name: user,
				ContextInfo: ContextInfo{
					Cluster: "minikube",
					User:    user,
				},
			},
		},
		CurrentContext: user,
		Users: []User{
			{
				Name: user,
				UserInfo: UserInfo{
					Token: token,
				},
			},
		},
	}
	b, err := yaml.Marshal(config)
	if err != nil {
		logrus.Error(err)
		return err
	} else {
		err = ioutil.WriteFile("kubeconfig.yaml", b, 0644)
		if err != nil {
			logrus.Error(err)
		}
		return err
	}
}
