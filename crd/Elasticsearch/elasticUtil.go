package Elasticsearch

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type ElasticsearchApi struct {
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

func (api *ElasticsearchApi) Get(namespace, name string) (*Elasticsearch, error) {
	result := Elasticsearch{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource("elasticsearches").
		Name(name).
		VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

func (api *ElasticsearchApi) List(namespace string) (*ElasticsearchList, error) {
	result := ElasticsearchList{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource("elasticsearches").
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

func (api *ElasticsearchApi) Delete(namespace, name string) error {
	err := api.Client.Delete().
		Namespace(namespace).
		Resource("elasticsearches").
		VersionedParams(&metav1.DeleteOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).Error()
	return err
}

func (api *ElasticsearchApi) Update(namespace, name string, other *Elasticsearch) error {
	err := api.Client.Put().
		Namespace(namespace).
		Resource("elasticsearches").
		Name(name).
		VersionedParams(&metav1.UpdateOptions{}, scheme.ParameterCodec).
		Body(other).
		Do(context.TODO()).
		Error()
	return err
}

// following methods maybe needed
func (api *ElasticsearchApi) Watch(opts metav1.ListOptions, namespace string) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return api.Client.Get().
		Namespace(namespace).
		Resource("elasticsearches").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(context.TODO())
}

func (api *ElasticsearchApi) WatchResources(namespace string) cache.Store {
	esStore, esController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return api.List(namespace)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return api.Watch(lo, namespace)
			},
		},
		&Elasticsearch{},
		40*time.Second,
		cache.ResourceEventHandlerFuncs{},
	)

	go esController.Run(wait.NeverStop)
	return esStore
}
