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

func (api *ElasticsearchApi) Get(opts metav1.GetOptions, namespace, name string) (*Elasticsearch, error) {
	result := Elasticsearch{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource("elasticsearches").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

func (api *ElasticsearchApi) List(opts metav1.ListOptions, namespace string) (*ElasticsearchList, error) {
	result := ElasticsearchList{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource("elasticsearches").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

func (api *ElasticsearchApi) Delete(opts metav1.DeleteOptions, namespace, name string) error {
	err := api.Client.Delete().
		Namespace(namespace).
		Resource("elasticsearches").
		VersionedParams(&opts, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).Error()
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
				return api.List(lo, namespace)
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
