package Postgresql

import (
	"context"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type PostgresqlApi struct {
	Client *rest.RESTClient
}

func NewConfigFor(kubeconfig string) (*rest.RESTClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logrus.WithError(err).Fatal("could not get config")
		return nil, err
	}
	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: GroupName, Version: GroupVersion}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	return rest.UnversionedRESTClientFor(&crdConfig)
}

// methods of PostgresqlApi struct ...

// clusters
func (api *PostgresqlApi) GetCluster(namespace, name string) (*Pgcluster, error) {
	result := Pgcluster{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgclusterResourcePlural).
		Name(name).
		VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

func (api *PostgresqlApi) ListCluster(namespace string) (*PgclusterList, error) {
	result := PgclusterList{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgclusterResourcePlural).
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

func (api *PostgresqlApi) DeleteCluster(namespace, name string) error {
	err := api.Client.Delete().
		Namespace(namespace).
		Resource(PgclusterResourcePlural).
		VersionedParams(&metav1.DeleteOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).Error()
	return err
}

// policies
func (api *PostgresqlApi) GetPolicy(namespace, name string) (*Pgpolicy, error) {
	result := Pgpolicy{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgpolicyResourcePlural).
		Name(name).
		VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

func (api *PostgresqlApi) ListPolicy(namespace string) (*PgpolicyList, error) {
	result := PgpolicyList{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgpolicyResourcePlural).
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

func (api *PostgresqlApi) DeletePolicy(namespace, name string) error {
	err := api.Client.Delete().
		Namespace(namespace).
		Resource(PgpolicyResourcePlural).
		VersionedParams(&metav1.DeleteOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).Error()
	return err
}

// replicas
func (api *PostgresqlApi) GetReplica(namespace, name string) (*Pgreplica, error) {
	result := Pgreplica{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgreplicaResourcePlural).
		Name(name).
		VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

func (api *PostgresqlApi) ListReplica(namespace string) (*PgreplicaList, error) {
	result := PgreplicaList{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgreplicaResourcePlural).
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

func (api *PostgresqlApi) DeleteReplica(namespace, name string) error {
	err := api.Client.Delete().
		Namespace(namespace).
		Resource(PgreplicaResourcePlural).
		VersionedParams(&metav1.DeleteOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).Error()
	return err
}

// tasks
func (api *PostgresqlApi) GetTask(namespace, name string) (*Pgtask, error) {
	result := Pgtask{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgtaskResourcePlural).
		Name(name).
		VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

func (api *PostgresqlApi) ListTask(namespace string) (*PgtaskList, error) {
	result := PgtaskList{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgtaskResourcePlural).
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

func (api *PostgresqlApi) DeleteTask(namespace, name string) error {
	err := api.Client.Delete().
		Namespace(namespace).
		Resource(PgtaskResourcePlural).
		VersionedParams(&metav1.DeleteOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).Error()
	return err
}
