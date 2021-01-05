package main

import (
	"fmt"

	//"k8s.io/client-go/tools/clientcmd"

	"goclient/K8s"
	msgs "goclient/apiservermsgs"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {

	/*
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			fmt.Println("here")
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
	*/
	GenerateKubeconfigTest()
}

// test kibana api
func kibanaApiTest() {
	api := Api{Kubeconfig: "/home/tung/.kube/config"}
	k8sApi, err := api.K8sAPI(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	ns, name := "default", "mybu"

	comm, err := api.GetCommunication(k8sApi, "kb", ns, name)

	if err != nil {
		fmt.Println(err)
		return
	}

	secret, err := k8sApi.GetSecret(metav1.GetOptions{}, ns, "wibu-es-elastic-user")
	if err != nil {
		fmt.Println(err)
		return
	}
	username := "elastic"
	password := string(secret.Data[username])
	_ = []byte(`{
		"id": "marketing",
		"name": "Marketing",
		"description" : "This is the Marketing Space",
		"color": "#aabbcc",
		"initials": "MK",
		"disabledFeatures": ["updated"],
		"imageUrl": ""
	}`)
	resp, err := comm.GetFeatures(username, password, 30928)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}

// test elasticsearch api
func elasticApiTest() {
	api := Api{Kubeconfig: "/home/tung/.kube/config"}
	k8sApi, err := api.K8sAPI(true)
	if err != nil {
		fmt.Println(err)
		return
	}
	ns, name := "default", "wibu"
	comm, err := api.GetCommunication(k8sApi, "es", ns, name)

	if err != nil {
		fmt.Println(err)
		return
	}

	secret, err := k8sApi.GetSecret(metav1.GetOptions{}, ns, "wibu-es-elastic-user")
	if err != nil {
		fmt.Println(err)
		return
	}
	username := "elastic"
	password := string(secret.Data[username])

	resp, err := comm.GetConnection(username, password, 30927)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}

// test postgresql api
func postgresApiTest() {
	api := Api{Kubeconfig: "/home/tung/.kube/config"}
	k8sApi, err := api.K8sAPI(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	ns := "pgo"
	comm, err := api.GetCommunication(k8sApi, "pg", ns, "")
	if err != nil {
		fmt.Println(err)
		return
	}

	username, password := "admin", "examplepassword"
	req := msgs.ShowClusterRequest{
		Namespace:     "pgo",
		ClientVersion: "4.5.1",
		AllFlag:       true,
	}
	resp, err := comm.ShowClusters(username, password, 32508, req)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(*resp)
}

// test LoadBalancer service create
func lbsvcTest() {
	api := Api{Kubeconfig: "/home/tung/.kube/config"}
	k8sApi, err := api.K8sAPI(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	service, err := k8sApi.CreateLBService(metav1.CreateOptions{}, "pg", "pgo", "postgres-operator", 8443)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(*service)
}

// test kibana crd tracking api
func kibanaCRDTest() {
	api := Api{Kubeconfig: "/home/tung/.kube/config"}
	kbApi, err := api.KibanaAPI()
	if err != nil {
		fmt.Println(err)
		return
	}
	kb, err := kbApi.Get(metav1.GetOptions{}, "default", "mybu")

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(kb)
}

// test elasticsearch crd tracking api
func elasticCRDTest() {
	api := Api{Kubeconfig: "/home/tung/.kube/config"}
	esApi, err := api.ElasticsearchAPI()
	if err != nil {
		fmt.Println(err)
		return
	}
	list, err := esApi.List(metav1.ListOptions{}, "default")

	if err != nil {
		fmt.Println(err)
		return
	}

	if list.Items != nil {
		for _, item := range list.Items {
			fmt.Println(item)
		}
	}
}

func GenerateKubeconfigTest() {
	/*
		// kubernetes-dashboard service account with full permission (admin-cluster)
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6IkROekxPY2NiVjA4UTljeVhnM0tOenZPNU5RUVlkbUlvNDV1WmJuNG1KME0ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlcm5ldGVzLWRhc2hib2FyZCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWRhc2hib2FyZC10b2tlbi01dnB4diIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWRhc2hib2FyZCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjIyZTNkMDg5LWZhOWEtNDZkYy04ODhkLWIxYjAxMzIxMGJmNSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlcm5ldGVzLWRhc2hib2FyZDprdWJlcm5ldGVzLWRhc2hib2FyZCJ9.nr1GKrJkjSngEKSf5fcX2EPVoVn-m9O69eIdOBf5udJArhBf74J3r_HkTTQ9HYhGMUAZFWDBacf6rjadCTuu7VxBWxkOyWB_MocjecXNogDeVowy9NXSQQofM0Da8VqNhIA5Fne4cgJyOzBvbZL1K9yKsgCTiMwHz0m67L8a2twFlScdSNHvGI8K2qvAH0MZLYHlnN8HgzM0Nbmz4r-eR3qlHlwTwP3N13ep_C1M1DMGJBVZ4fz2ntR0r7AO-3lkj8oL5S3APlZ21aWwRaHMbxq6OXlX-49IqMBzj0m7GgAS4ZI2M_Xf3iI-SxRstgp6DDHABL69NVZy6TwSKCOqig"
		user := "kubernetes-dashboard"
	*/

	// my service account with permission to just namespace lol
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6IkROekxPY2NiVjA4UTljeVhnM0tOenZPNU5RUVlkbUlvNDV1WmJuNG1KME0ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJsb2wiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlY3JldC5uYW1lIjoiYnVpbGQtcm9ib3QtdG9rZW4tcWtiengiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiYnVpbGQtcm9ib3QiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIzZjhkOWQ4Ny04Njk2LTQ5NzYtOGM1MS0wY2FjODM0OTkwZTUiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6bG9sOmJ1aWxkLXJvYm90In0.DmfcKUqW7sH8OVHM8YJjerLbqj2HdCqo-_37Kpv-5_AIkr10Zp4dQDUyrC0VQ2C8P6dlh5lEEDNpHofS7W33DXpmb-jzqGF8ETQB_XQJDosor7hzm7UGsBCz7A2U7Fp8H-CFwFAMae0M_Lna1Rsz9F6587VXZJuUlDZqL7RVu2UvJuf4vksJ0Ht0leYL4N5dtR5kwySGvM9CUUpuJvEudIHbcSjDB1h55rOSPH4DlxZFjXtX028Xv1SBzdt9IOIkVJ5ZfXZRjiCHanIRKRrrCjIkEvWybKfNpVhRT_eQ0JshtSngJKHwhR72OHMtgZG4Kf-mcKDSiidytRy9vcok1Q"
	user := "build-robot"

	err := K8s.GenerateKubeconfig(token, user)
	if err != nil {
		fmt.Println(err)
		return
	}

	api := Api{Kubeconfig: "kubeconfig.yaml"}
	k8sApi, err := api.K8sAPI(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	// test whether can get service from namespace default
	list, err := k8sApi.PollServices("default")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, svc := range list.Items {
		fmt.Println(svc)
	}

}
