package main

import (
	"fmt"

	//"k8s.io/client-go/tools/clientcmd"

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
	InitRoute()
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
	data := msgs.DeleteClusterRequest{
		Clustername:   "ganmo",
		Namespace:     "pgo",
		ClientVersion: "4.5.1",
		DeleteBackups: true,
		DeleteData:    true,
	}
	resp, err := comm.DeleteClusters(username, password, 32508, data)

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

func kibanaCRDTest() {
	api := Api{Kubeconfig: "/home/tung/.kube/config"}
	kbApi, err := api.KibanaAPI()
	if err != nil {
		fmt.Println(err)
		return
	}
	list, err := kbApi.List(metav1.ListOptions{}, "default")

	if err != nil {
		fmt.Println(err)
		return
	}

	if list.Items != nil {
		for _, item := range list.Items {
			fmt.Println(item.Status)
		}
	}
}

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
