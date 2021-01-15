package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	comm "goclient/Communication"
	"goclient/K8s"
	"goclient/RestStruct/Postgres/structs"
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
	kibanaCRDTest()
}

// test kibana api
func kibanaApiTest() {
	api := comm.Api{Kubeconfig: "/home/tung/.kube/config"}
	k8sApi, err := api.K8sAPI()
	if err != nil {
		fmt.Println(err)
		return
	}

	ns, name := "default", "mybu"

	comm, err := api.GetCommunication("kb", ns, name)

	if err != nil {
		fmt.Println(err)
		return
	}

	secret, err := k8sApi.GetSecret(ns, "wibu-es-elastic-user")
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
	resp, err := comm.GetSpace(username, password, "marketing")

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}

// test elasticsearch api
func elasticApiTest() {
	api := comm.Api{Kubeconfig: "/home/tung/.kube/config"}
	k8sApi, err := api.K8sAPI()
	if err != nil {
		fmt.Println(err)
		return
	}
	ns, name := "default", "wibu"
	comm, err := api.GetCommunication("es", ns, name)

	if err != nil {
		fmt.Println(err)
		return
	}

	secret, err := k8sApi.GetSecret(ns, "wibu-es-elastic-user")
	if err != nil {
		fmt.Println(err)
		return
	}
	username := "elastic"
	password := string(secret.Data[username])

	resp, err := comm.GetConnection(username, password)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}

// test postgresql api
func postgresApiTest() {
	api := comm.Api{Kubeconfig: "/home/tung/.kube/config"}

	ns := "pgo"
	comm, err := api.GetCommunication("pg", ns, "postgres-operator")
	if err != nil {
		fmt.Println(err)
		return
	}

	username, password := "admin", "examplepassword"
	req := structs.ShowClusterRequest{
		Namespace:     "pgo",
		ClientVersion: "4.5.1",
		AllFlag:       true,
	}
	resp, err := comm.ShowClusters(username, password, req)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(*resp)
}

// test kibana crd tracking api
func kibanaCRDTest() {
	api := comm.Api{Kubeconfig: "/home/tung/.kube/config"}
	kbApi, err := api.KibanaAPI()
	if err != nil {
		fmt.Println(err)
		return
	}
	ns, name := "default", "mybu"
	kb, err := kbApi.Get(ns, name)

	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := json.Marshal(kb)
	if err != nil {
		fmt.Println(err)
		return
	}
	ioutil.WriteFile("./files/tmp.json", b, 0644)
}

// test elasticsearch crd tracking api
func elasticCRDTest() {
	api := comm.Api{Kubeconfig: "/home/tung/.kube/config"}
	esApi, err := api.ElasticsearchAPI()
	if err != nil {
		fmt.Println(err)
		return
	}
	es, err := esApi.Get("default", "wibu")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(es)
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

	api := comm.Api{Kubeconfig: "kubeconfig.yaml"}
	k8sApi, err := api.K8sAPI()
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

func ScaleTest() {
	api := comm.Api{Kubeconfig: "/home/tung/.kube/config"}
	k8sApi, _ := api.K8sAPI()
	scale, _ := k8sApi.GetStatefulSetScale("default", "wibu-es-ganmo")
	fmt.Println(scale)
}

func IngressTest() {
	api := comm.Api{Kubeconfig: "/home/tung/.kube/config"}
	k8sApi, _ := api.K8sAPI()
	ns := "default"
	err := k8sApi.AddServiceToIngress(ns, ns+"-ingress", "wibu-es-http", "ganmo.com", 9200)
	if err != nil {
		fmt.Println(err)
	}
	cert, _ := ioutil.ReadFile("./serverCerts/server.crt")
	key, _ := ioutil.ReadFile("./serverCerts/server.key")
	_, err = k8sApi.CreateSecret(ns, ns+"-ingress-sercret", cert, key)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateKibanaTest() {
	api := comm.Api{Kubeconfig: "/home/tung/.kube/config"}
	k8sApi, _ := api.K8sAPI()
	_, err := k8sApi.ApplyFile("kibana.yaml", "apply")
	if err != nil {
		fmt.Println(err)
		return
	}
	name, ns, kind := "mybu", "default", "kb"
	svc, err := k8sApi.GetService(ns, name)
	for svc == nil {
		svc, err = k8sApi.GetService(ns, name)
		time.Sleep(time.Second * 4)
	}
	err = k8sApi.AddServiceToIngress(ns, ns+"-ingress", name+"-"+kind+"-http", "ganmo.com", 5601)
	if err != nil {
		fmt.Println(err)
		return
	}
}
