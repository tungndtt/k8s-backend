package K8s

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type KubeConfig struct {
	ApiVersion     string    `yaml:"apiVersion"`
	Kind           string    `yaml:"kind"`
	Clusters       []Cluster `yaml:"clusters"`
	Contexts       []Context `yaml:"contexts"`
	CurrentContext string    `yaml:"current-context"`
	Users          []User    `yaml:"users"`
}

type Cluster struct {
	ClusterInfo ClusterInfo `yaml:"cluster"`
	Name        string      `yaml:"name"`
}

type ClusterInfo struct {
	Cert      string `yaml:"certificate-authority"`
	ServerUrl string `yaml:"server"`
}

type Context struct {
	ContextInfo ContextInfo `yaml:"context"`
	Name        string      `yaml:"name"`
}

type ContextInfo struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type User struct {
	Name     string   `yaml:"name"`
	UserInfo UserInfo `yaml:"user"`
}

type UserInfo struct {
	Token string `yaml:"token"`
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
					ServerUrl: "https://192.168.99.112:8443",
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
		return err
	} else {
		return ioutil.WriteFile("kubeconfig.yaml", b, 0644)
	}
}
