package Communication

import (
	"goclient/Elasticsearch"
	"goclient/K8s"
	"goclient/Kibana"
	"goclient/Postgresql"
)

type Api struct {
	Kubeconfig string
}

func (api *Api) K8sAUTH() (*K8s.K8sAuth, error) {
	client, err := K8s.GetRbacClient(api.Kubeconfig)
	if err != nil {
		return nil, err
	} else {
		return &K8s.K8sAuth{RbacClient: client}, nil
	}
}

func (api *Api) K8sAPI(outside bool) (*K8s.K8sApi, error) {
	clientset, err1, dif, err2, config := K8s.GetClientSet(outside, api.Kubeconfig)
	if err1 != nil {
		return nil, err1
	} else if err2 != nil {
		return nil, err2
	} else {
		return &K8s.K8sApi{
			ClientSet: clientset,
			Dif:       dif,
			Config:    config,
		}, nil
	}
}

func (api *Api) KibanaAPI() (*Kibana.KibanaApi, error) {
	client, err := Kibana.NewConfigFor(api.Kubeconfig)
	if err != nil {
		return nil, err
	} else {
		return &Kibana.KibanaApi{
			Client: client,
		}, nil
	}
}

func (api *Api) ElasticsearchAPI() (*Elasticsearch.ElasticsearchApi, error) {
	client, err := Elasticsearch.NewConfigFor(api.Kubeconfig)
	if err != nil {
		return nil, err
	} else {
		return &Elasticsearch.ElasticsearchApi{
			Client: client,
		}, nil
	}
}

func (api *Api) PostgresqlAPI() (*Postgresql.PostgresqlApi, error) {
	client, err := Postgresql.NewConfigFor(api.Kubeconfig)
	if err != nil {
		return nil, err
	} else {
		return &Postgresql.PostgresqlApi{
			Client: client,
		}, nil
	}
}

func (api *Api) GetCommunication(kind, namespace, name string) (*Comm, error) {
	return getCommunication(api.Kubeconfig, kind, namespace, name)
}
