package K8s

import (
	"context"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"

	"bytes"
	"net/http"
	"os"
	"time"

	"encoding/json"

	appV1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

type K8sApi struct {
	ClientSet *kubernetes.Clientset
	Dif       dynamic.Interface
	Config    *rest.Config
}

func GetClientSet(outside bool, kubeconfig string) (*kubernetes.Clientset, error, dynamic.Interface, error, *rest.Config) {
	var err error
	var Config *rest.Config
	if !outside && kubeconfig == "" {
		logrus.Info("Using inside Cluster")
		Config, err = rest.InClusterConfig()
	} else {
		logrus.Info("Using outside cluster")
		Config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		return nil, err, nil, err, nil
	}

	clientSet, err1 := kubernetes.NewForConfig(Config)
	mDynamic, err2 := dynamic.NewForConfig(Config)
	if err1 != nil || err2 != nil {
		return nil, err1, nil, err2, Config
	}
	return clientSet, nil, mDynamic, nil, Config
}

// pass byte[] in param instead of filePath (because we dont save the file)
func (api *K8sApi) ApplyFile(filePath, opCode string) (*unstructured.Unstructured, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(b), 100)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			return nil, err
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if obj == nil {
			return nil, err
		}
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(api.ClientSet.Discovery())
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace("default")
			}
			dri = api.Dif.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = api.Dif.Resource(mapping.Resource)
		}

		if opCode == "create" {
			unstructured, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{
				FieldManager: "field-manager",
			})
			if err != nil {
				logrus.Error(err)
			}
			return unstructured, err
		} else if opCode == "apply" {
			data, _ := json.Marshal(obj)
			force := true
			unstructured, err := dri.Patch(context.Background(), unstructuredObj.GetName(),
				types.ApplyPatchType, data, metav1.PatchOptions{
					FieldManager: "field-manager",
					Force:        &force,
				})
			if err != nil {
				logrus.Error(err)
			} else {
				logrus.Info((unstructured.Object["spec"].(map[string]interface{}))["version"])
			}
			return unstructured, err
		} else if opCode == "delete" {
			err = dri.Delete(context.Background(), unstructuredObj.GetName(), metav1.DeleteOptions{})
			if err != nil {
				logrus.Error(err)
			} else {
				logrus.Info("Delete successfully")
			}
			return nil, err
		}
	}
}

func (api *K8sApi) GetPod(opts metav1.GetOptions, namespace, podname string) (*v1.Pod, error) {
	return api.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podname, opts)
}

func (api *K8sApi) GetService(opts metav1.GetOptions, namespace, servicename string) (*v1.Service, error) {
	return api.ClientSet.CoreV1().Services(namespace).Get(context.TODO(), servicename, opts)
}

func (api *K8sApi) GetServiceAccount(opts metav1.GetOptions, namespace, name string) (*v1.ServiceAccount, error) {
	return api.ClientSet.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), name, opts)
}

func (api *K8sApi) GetReplicaset(opts metav1.GetOptions, namespace, name string) (*appV1.ReplicaSet, error) {
	return api.ClientSet.AppsV1().ReplicaSets(namespace).Get(context.TODO(), name, opts)
}

func (api *K8sApi) GetStatefulset(opts metav1.GetOptions, namespace, name string) (*appV1.StatefulSet, error) {
	return api.ClientSet.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, opts)
}

func (api *K8sApi) GetSecret(opts metav1.GetOptions, namespace, name string) (*v1.Secret, error) {
	return api.ClientSet.CoreV1().Secrets(namespace).Get(context.TODO(), name, opts)
}

func (api *K8sApi) GetDeploymentScale(opts metav1.GetOptions, namespace, name string) (*autoscalingv1.Scale, error) {
	return api.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), name, opts)
}

func (api *K8sApi) GetStatefulSetScale(opts metav1.GetOptions, namespace, name string) (*autoscalingv1.Scale, error) {
	return api.ClientSet.AppsV1().StatefulSets(namespace).GetScale(context.TODO(), name, opts)
}

func (api *K8sApi) CreateServiceAccount(opts metav1.CreateOptions, namespace, name string) (*v1.ServiceAccount, error) {
	new := v1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	return api.ClientSet.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), &new, opts)
}

func (api *K8sApi) CreateLBService(opts metav1.CreateOptions, kind, namespace, name, uid string, port int32) (*v1.Service, error) {
	var selector map[string]string
	var ownerReferences []metav1.OwnerReference
	if kind == "es" {
		selector = map[string]string{
			"elasticsearch.k8s.elastic.co/cluster-name": name,
			"common.k8s.elastic.co/type":                "elasticsearch",
		}
		controller, blockOwnerDeletion := true, true
		ownerReferences = []metav1.OwnerReference{
			{
				APIVersion:         "elasticsearch.k8s.elastic.co/v1",
				Kind:               "Elasticsearch",
				Name:               name,
				UID:                types.UID(uid),
				Controller:         &controller,
				BlockOwnerDeletion: &blockOwnerDeletion,
			},
		}
	} else if kind == "kb" {
		selector = map[string]string{
			"common.k8s.elastic.co/type": "kibana",
			"kibana.k8s.elastic.co/name": name,
		}
		controller, blockOwnerDeletion := true, true
		ownerReferences = []metav1.OwnerReference{
			{
				APIVersion:         "kibana.k8s.elastic.co/v1",
				Kind:               "Kibana",
				Name:               name,
				UID:                types.UID(uid),
				Controller:         &controller,
				BlockOwnerDeletion: &blockOwnerDeletion,
			},
		}
	} else if kind == "pg" {
		selector = map[string]string{
			"name":   name,
			"vendor": "crunchydata",
		}
	}
	new := v1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            name + "-service",
			Namespace:       namespace,
			GenerateName:    name + "-service",
			OwnerReferences: ownerReferences,
		},
		Spec: v1.ServiceSpec{
			Selector: selector,
			Type:     "LoadBalancer",
			Ports: []v1.ServicePort{
				{
					Name: name,
					Port: port,
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: port,
					},
					Protocol: "TCP",
				},
			},
		},
	}
	return api.ClientSet.CoreV1().Services(namespace).Create(context.TODO(), &new, opts)
}

func (api *K8sApi) PollPods(namespace string) (*v1.PodList, error) {
	return api.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (api *K8sApi) PollServices(namespace string) (*v1.ServiceList, error) {
	return api.ClientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (api *K8sApi) PollReplicaSets(opts metav1.ListOptions, namespace string) (*appV1.ReplicaSetList, error) {
	return api.ClientSet.AppsV1().ReplicaSets(namespace).List(context.TODO(), opts)
}

func (api *K8sApi) PollStatefulset(opts metav1.ListOptions, namespace string) (*appV1.StatefulSetList, error) {
	return api.ClientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), opts)
}

func (api *K8sApi) PollServiceAccounts(namespace string) (*v1.ServiceAccountList, error) {
	return api.ClientSet.CoreV1().ServiceAccounts(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (api *K8sApi) DeleteServiceAccount(opts metav1.DeleteOptions, namespace, name string) error {
	return api.ClientSet.CoreV1().ServiceAccounts(namespace).Delete(context.TODO(), name, opts)
}

func (api *K8sApi) DeleteService(opts metav1.DeleteOptions, namespace, name string) error {
	return api.ClientSet.CoreV1().Services(namespace).Delete(context.TODO(), name, opts)
}

func (api *K8sApi) UpdateScaleDeployment(opts metav1.UpdateOptions, namespace, name string, num int32) (*autoscalingv1.Scale, error) {
	scale, err := api.GetDeploymentScale(metav1.GetOptions{}, namespace, name)
	if err != nil {
		return nil, err
	}
	scale.Spec.Replicas = num
	return api.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), name, scale, opts)
}

func (api *K8sApi) UpdateScaleStatefulSet(opts metav1.UpdateOptions, namespace, name string, num int32) (*autoscalingv1.Scale, error) {
	scale, err := api.GetStatefulSetScale(metav1.GetOptions{}, namespace, name)
	if err != nil {
		return nil, err
	}
	scale.Spec.Replicas = num
	return api.ClientSet.AppsV1().StatefulSets(namespace).UpdateScale(context.TODO(), name, scale, opts)
}

// following methods maybe needed
func (api *K8sApi) OpenPortFowarding(namespace, name string) string {
	reqURL := api.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(namespace).
		Name(name).
		SubResource("portforward").URL()

	logrus.Infof("got url: %s", reqURL)
	transport, upgrader, _ := spdy.RoundTripperFor(api.Config)

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, reqURL)

	stopChan, readyChan := make(<-chan struct{}, 1), make(chan struct{}, 1)
	out, errOut := os.Stdout, os.Stdout
	fw, _ := portforward.New(dialer, []string{"9200:5601"}, stopChan, readyChan, out, errOut)
	fw.ForwardPorts()
	return reqURL.String()
}

func (api *K8sApi) WatchPods() {
	logrus.Info("watching")
	factory := informers.NewSharedInformerFactory(api.ClientSet, time.Second*30)
	informer := factory.Core().V1().Pods().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{

		UpdateFunc: func(oldObj, newObj interface{}) {
			logrus.Infof("Updating ...", oldObj, newObj)
		},
		AddFunc: func(obj interface{}) {
			logrus.Infof("Adding ...", obj)
		},
		DeleteFunc: func(obj interface{}) {
			logrus.Infof("Deleting ...", obj)
		},
	})
	stopper := make(chan struct{})
	go informer.Run(stopper)
}
