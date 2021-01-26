package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	api "goclient/Api"
	pgStructs "goclient/RestStruct/Postgres/structs"
	"goclient/files"
	_ "goclient/internal"
	esSvc "goclient/internal/elasticsearch"
	kbSvc "goclient/internal/kibana"
	pgSvc "goclient/internal/postgres"
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
	elasticApiTest()
}

// test kibana api
func kibanaApiTest() {
	api := api.Api{Kubeconfig: "/home/tung/.kube/config"}
	k8sApi, _ := api.K8sAPI()
	kbApi, _ := api.KibanaAPI()
	ns, name := "default", "mybu"
	svc, _ := kbSvc.GetKibanaService(k8sApi, kbApi, ns, name)

	_ = []byte(`{
		"id": "marketing",
		"name": "Marketing",
		"description" : "This is the Marketing Space",
		"color": "#aabbcc",
		"initials": "MK",
		"disabledFeatures": ["updated"],
		"imageUrl": ""
	}`)
	resp, err := svc.ExecuteCustomAction("ACTION_GET_FEATURES", []byte(`Get all features`))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}

// test elasticsearch api
func elasticApiTest() {
	api := api.Api{Kubeconfig: "/home/tung/.kube/config"}
	k8sApi, _ := api.K8sAPI()
	esApi, _ := api.ElasticsearchAPI()
	ns, name := "default", "wibu"
	svc, _ := esSvc.GetElasticsearchService(k8sApi, esApi, ns, name)

	resp, err := svc.ExecuteCustomAction("ACTION_GET_CONNECTION", []byte("test"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}

// test postgresql api
func postgresApiTest() {
	api := api.Api{Kubeconfig: "/home/tung/.kube/config"}
	k8sApi, _ := api.K8sAPI()
	pgApi, _ := api.PostgresqlAPI()
	ns, name := "pgo", "postgres-operator"
	_, err := pgSvc.GetPostgresService(k8sApi, pgApi, ns, name)

	if err != nil {
		fmt.Println(err)
		return
	}

	_ = pgStructs.ShowClusterRequest{
		Namespace:     "pgo",
		ClientVersion: "4.5.1",
		AllFlag:       true,
	}

	//service.ExecuteCustomAction()
	/*
		Namespace:     "default",
		ClientVersion: "4.5.1",
		Name:          "a_test_cluster",
		Username:      "tung",
		Password:      "tung",
	*/

	/*
		comm, err := svc.

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(*resp)
	*/
}

// test kibana crd tracking api
func kibanaCRDTest() {
	api := api.Api{Kubeconfig: "/home/tung/.kube/config"}
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
	api := api.Api{Kubeconfig: "/home/tung/.kube/config"}
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

	err := files.GenerateKubeconfig(token, user)
	if err != nil {
		fmt.Println(err)
		return
	}

	api := api.Api{Kubeconfig: "./files/kubeconfig.yaml"}
	k8sApi, err := api.K8sAPI()
	if err != nil {
		fmt.Println(err)
		return
	}

	// test whether can get service from namespace default
	list, err := k8sApi.PollServices("lol")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, svc := range list.Items {
		fmt.Println(svc)
	}
}

func ScaleTest() {
	api := api.Api{Kubeconfig: "/home/tung/.kube/config"}
	k8sApi, _ := api.K8sAPI()
	scale, _ := k8sApi.GetStatefulSetScale("default", "wibu-es-ganmo")
	fmt.Println(scale)
}

func IngressTest() {
	api := api.Api{Kubeconfig: "/home/tung/.kube/config"}
	k8sApi, _ := api.K8sAPI()
	ns, name, kind := "pgo", "postgres-operator", "pg"
	err := k8sApi.AddServiceToIngress(ns, ns+"-ingress", fmt.Sprintf("%s-%s-http", name, kind), "ganmo.com", 8443)
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
	api := api.Api{Kubeconfig: "/home/tung/.kube/config"}
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
