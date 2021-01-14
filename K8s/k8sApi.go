package K8s

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	rbac "k8s.io/client-go/kubernetes/typed/rbac/v1"
)

type K8sApi struct {
	ClientSet     *kubernetes.Clientset
	Dif           dynamic.Interface
	V1beta1Client *v1beta1.ExtensionsV1beta1Client
	RbacClient    *rbac.RbacV1Client
}
