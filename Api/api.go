package Api

import (
	"goclient/K8s"
	"goclient/crd/Elasticsearch"
	"goclient/crd/Kibana"
	"goclient/crd/Postgresql"
)

type Api struct {
	Kubeconfig string
}

func (api *Api) K8sAPI() (*K8s.K8sApi, error) {
	config, err := K8s.GenerateConfig(api.Kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err1, dif, err2 := K8s.GetClientSet(config)
	if err1 != nil {
		return nil, err1
	} else if err2 != nil {
		return nil, err2
	} else {
		rbacClient, err := K8s.GetRbacClient(config)
		if err != nil {
			return nil, err
		}
		v1beta1Client, err := K8s.GetV1Beta1Client(config)
		if err != nil {
			return nil, err
		}
		return &K8s.K8sApi{
			ClientSet:     clientset,
			Dif:           dif,
			RbacClient:    rbacClient,
			V1beta1Client: v1beta1Client,
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
