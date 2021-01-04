package K8s

import (
	"context"

	"github.com/sirupsen/logrus"

	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbac "k8s.io/client-go/kubernetes/typed/rbac/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sAuth struct {
	RbacClient *rbac.RbacV1Client
}

func GetRbacClient(kubeconfig string) (*rbac.RbacV1Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return rbac.NewForConfig(config)
}

func (auth *K8sAuth) CreateRoleBinding(opts metav1.CreateOptions, namespace, name, rolename string, account_namespace, account_name []string) (*rbacV1.RoleBinding, error) {
	subjects := []rbacV1.Subject{}
	for i := 0; i < len(account_name); i++ {
		subjects = append(
			subjects,
			rbacV1.Subject{
				Kind:      "ServiceAccount",
				Name:      account_name[i],
				Namespace: account_namespace[i],
			},
		)
	}
	new := rbacV1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		RoleRef: rbacV1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     rolename,
		},
		Subjects: subjects,
	}
	return auth.RbacClient.RoleBindings(namespace).Create(context.TODO(), &new, opts)
}

func (auth *K8sAuth) CreateRole(opts metav1.CreateOptions, namespace, name string, apiGroups, resources, verbs [][]string) (*rbacV1.Role, error) {
	rules := []rbacV1.PolicyRule{}
	for i := 0; i < len(verbs); i++ {
		rules = append(
			rules,
			rbacV1.PolicyRule{
				APIGroups: apiGroups[i],
				Resources: resources[i],
				Verbs:     verbs[i],
			},
		)
	}
	new := rbacV1.Role{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Role",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Rules: rules,
	}
	return auth.RbacClient.Roles(namespace).Create(context.TODO(), &new, opts)
}

func (auth *K8sAuth) GetRoleBinding(opts metav1.GetOptions, namespace, name string) (*rbacV1.RoleBinding, error) {
	return auth.RbacClient.RoleBindings(namespace).Get(context.TODO(), name, opts)
}

func (auth *K8sAuth) GetRole(opts metav1.GetOptions, namespace, name string) (*rbacV1.Role, error) {
	return auth.RbacClient.Roles(namespace).Get(context.TODO(), name, opts)
}

func (auth *K8sAuth) AddUserTo(opts metav1.UpdateOptions, old *rbacV1.RoleBinding, namespace, name, apiGroup string) error {
	new := rbacV1.RoleBinding{
		TypeMeta: old.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:      old.ObjectMeta.Name,
			Namespace: old.ObjectMeta.Namespace,
			SelfLink:  old.ObjectMeta.SelfLink,
		},
		Subjects: old.Subjects,
		RoleRef:  old.RoleRef,
	}
	new.Subjects = append(new.Subjects, rbacV1.Subject{Kind: "ServiceAccount", APIGroup: apiGroup, Name: name, Namespace: namespace})
	_, err := auth.RbacClient.RoleBindings(namespace).Update(context.TODO(), &new, opts)
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func (auth *K8sAuth) DeleteUserFrom(old *rbacV1.RoleBinding, name, namespace, kind string) error {
	newSubjects := old.Subjects
	for i, sub := range newSubjects {
		if sub.Kind == kind && sub.Name == name && sub.Namespace == namespace {
			n := len(newSubjects)
			newSubjects[i], newSubjects[n-1] = newSubjects[n-1], newSubjects[i]
			newSubjects = newSubjects[:n-1]
			break
		}
	}
	new := rbacV1.RoleBinding{
		TypeMeta: old.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:      old.ObjectMeta.Name,
			Namespace: old.ObjectMeta.Namespace,
			SelfLink:  old.ObjectMeta.SelfLink,
		},
		Subjects: newSubjects,
		RoleRef:  old.RoleRef,
	}
	_, err := auth.RbacClient.RoleBindings(namespace).Update(context.TODO(), &new, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func (auth *K8sAuth) DeleteRole(opts metav1.DeleteOptions, namespace, name string) error {
	return auth.RbacClient.Roles(namespace).Delete(context.TODO(), name, opts)
}

func (auth *K8sAuth) DeleteRoleBinding(opts metav1.DeleteOptions, namespace, name string) error {
	return auth.RbacClient.RoleBindings(namespace).Delete(context.TODO(), name, opts)
}

func (auth *K8sAuth) PollRoles(opts metav1.ListOptions, namespace, name string) (*rbacV1.RoleList, error) {
	return auth.RbacClient.Roles(namespace).List(context.TODO(), opts)
}

func (auth *K8sAuth) PollRoleBindings(opts metav1.ListOptions, namespace, name string) (*rbacV1.RoleBindingList, error) {
	return auth.RbacClient.RoleBindings(namespace).List(context.TODO(), opts)
}
