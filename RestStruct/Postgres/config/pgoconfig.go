package config

/*
Copyright 2018 - 2020 Crunchy Data Solutions, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"

	crv1 "goclient/crd/Postgresql"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

const CustomConfigMapName = "pgo-config"
const DefaultConfigsPath = "/default-pgo-config/"
const CustomConfigsPath = "/pgo-config/"

var PgoDefaultServiceAccountTemplate *template.Template

const PGODefaultServiceAccountPath = "pgo-default-sa.json"

var PgoTargetRoleBindingTemplate *template.Template

const PGOTargetRoleBindingPath = "pgo-target-role-binding.json"

var PgoBackrestServiceAccountTemplate *template.Template

const PGOBackrestServiceAccountPath = "pgo-backrest-sa.json"

var PgoTargetServiceAccountTemplate *template.Template

const PGOTargetServiceAccountPath = "pgo-target-sa.json"

var PgoBackrestRoleTemplate *template.Template

const PGOBackrestRolePath = "pgo-backrest-role.json"

var PgoBackrestRoleBindingTemplate *template.Template

const PGOBackrestRoleBindingPath = "pgo-backrest-role-binding.json"

var PgoTargetRoleTemplate *template.Template

const PGOTargetRolePath = "pgo-target-role.json"

var PgoPgServiceAccountTemplate *template.Template

const PGOPgServiceAccountPath = "pgo-pg-sa.json"

var PgoPgRoleTemplate *template.Template

const PGOPgRolePath = "pgo-pg-role.json"

var PgoPgRoleBindingTemplate *template.Template

const PGOPgRoleBindingPath = "pgo-pg-role-binding.json"

var PolicyJobTemplate *template.Template

const policyJobTemplatePath = "pgo.sqlrunner-template.json"

var PVCTemplate *template.Template

const pvcPath = "pvc.json"

var ContainerResourcesTemplate *template.Template

const containerResourcesTemplatePath = "container-resources.json"

var AffinityTemplate *template.Template

const affinityTemplatePath = "affinity.json"

var PodAntiAffinityTemplate *template.Template

const podAntiAffinityTemplatePath = "pod-anti-affinity.json"

var PgoBackrestRepoServiceTemplate *template.Template

const pgoBackrestRepoServiceTemplatePath = "pgo-backrest-repo-service-template.json"

var PgoBackrestRepoTemplate *template.Template

const pgoBackrestRepoTemplatePath = "pgo-backrest-repo-template.json"

var PgmonitorEnvVarsTemplate *template.Template

const pgmonitorEnvVarsPath = "pgmonitor-env-vars.json"

var PgbackrestEnvVarsTemplate *template.Template

const pgbackrestEnvVarsPath = "pgbackrest-env-vars.json"

var PgbackrestS3EnvVarsTemplate *template.Template

const pgbackrestS3EnvVarsPath = "pgbackrest-s3-env-vars.json"

var PgAdminTemplate *template.Template

const pgAdminTemplatePath = "pgadmin-template.json"

var PgAdminServiceTemplate *template.Template

const pgAdminServiceTemplatePath = "pgadmin-service-template.json"

var PgbouncerTemplate *template.Template

const pgbouncerTemplatePath = "pgbouncer-template.json"

var PgbouncerConfTemplate *template.Template

const pgbouncerConfTemplatePath = "pgbouncer.ini"

var PgbouncerUsersTemplate *template.Template

const pgbouncerUsersTemplatePath = "users.txt"

var PgbouncerHBATemplate *template.Template

const pgbouncerHBATemplatePath = "pgbouncer_hba.conf"

var ServiceTemplate *template.Template

const serviceTemplatePath = "cluster-service.json"

var RmdatajobTemplate *template.Template

const rmdatajobPath = "rmdata-job.json"

var BackrestjobTemplate *template.Template

const backrestjobPath = "backrest-job.json"

var BackrestRestorejobTemplate *template.Template

const backrestRestorejobPath = "backrest-restore-job.json"

var PgDumpBackupJobTemplate *template.Template

const pgDumpBackupJobPath = "pgdump-job.json"

var PgRestoreJobTemplate *template.Template

const pgRestoreJobPath = "pgrestore-job.json"

var PVCMatchLabelsTemplate *template.Template

const pvcMatchLabelsPath = "pvc-matchlabels.json"

var PVCStorageClassTemplate *template.Template

const pvcSCPath = "pvc-storageclass.json"

var ExporterTemplate *template.Template

const exporterTemplatePath = "exporter.json"

var BadgerTemplate *template.Template

const badgerTemplatePath = "pgbadger.json"

var DeploymentTemplate *template.Template

const deploymentTemplatePath = "cluster-deployment.json"

var BootstrapTemplate *template.Template

const bootstrapTemplatePath = "cluster-bootstrap-job.json"

type ClusterStruct struct {
	CCPImagePrefix                 string
	CCPImageTag                    string
	Policies                       string
	Metrics                        bool
	Badger                         bool
	Port                           string
	PGBadgerPort                   string
	ExporterPort                   string
	User                           string
	Database                       string
	PasswordAgeDays                string
	PasswordLength                 string
	Replicas                       string
	ServiceType                    string
	BackrestPort                   int
	BackrestS3Bucket               string
	BackrestS3Endpoint             string
	BackrestS3Region               string
	BackrestS3URIStyle             string
	BackrestS3VerifyTLS            string
	DisableAutofail                bool
	PgmonitorPassword              string
	EnableCrunchyadm               bool
	DisableReplicaStartFailReinit  bool
	PodAntiAffinity                string
	PodAntiAffinityPgBackRest      string
	PodAntiAffinityPgBouncer       string
	SyncReplication                bool
	DefaultInstanceResourceMemory  resource.Quantity `json:"DefaultInstanceMemory"`
	DefaultBackrestResourceMemory  resource.Quantity `json:"DefaultBackrestMemory"`
	DefaultPgBouncerResourceMemory resource.Quantity `json:"DefaultPgBouncerMemory"`
	DefaultExporterResourceMemory  resource.Quantity `json:"DefaultExporterMemory"`
	DisableFSGroup                 bool
}

type StorageStruct struct {
	AccessMode         string
	Size               string
	StorageType        string
	StorageClass       string
	SupplementalGroups string
	MatchLabels        string
}

// PgoStruct defines various configuration settings for the PostgreSQL Operator
type PgoStruct struct {
	Audit                          bool
	ConfigMapWorkerCount           *int
	ControllerGroupRefreshInterval *int
	DisableReconcileRBAC           bool
	NamespaceRefreshInterval       *int
	NamespaceWorkerCount           *int
	PGClusterWorkerCount           *int
	PGOImagePrefix                 string
	PGOImageTag                    string
	PGReplicaWorkerCount           *int
	PGTaskWorkerCount              *int
}

type PgoConfig struct {
	BasicAuth       string
	Cluster         ClusterStruct
	Pgo             PgoStruct
	PrimaryStorage  string
	WALStorage      string
	BackupStorage   string
	ReplicaStorage  string
	BackrestStorage string
	Storage         map[string]StorageStruct
}

const DEFAULT_SERVICE_TYPE = "ClusterIP"
const LOAD_BALANCER_SERVICE_TYPE = "LoadBalancer"
const NODEPORT_SERVICE_TYPE = "NodePort"
const CONFIG_PATH = "pgo.yaml"

var log_statement_values = []string{"ddl", "none", "mod", "all"}

const DEFAULT_BACKREST_PORT = 2022
const DEFAULT_PGADMIN_PORT = "5050"
const DEFAULT_PGBADGER_PORT = "10000"
const DEFAULT_EXPORTER_PORT = "9187"
const DEFAULT_POSTGRES_PORT = "5432"
const DEFAULT_PATRONI_PORT = "8009"

func (c *PgoConfig) Validate() error {
	var err error
	errPrefix := "Error in pgoconfig: check pgo.yaml: "

	if c.Cluster.BackrestPort == 0 {
		c.Cluster.BackrestPort = DEFAULT_BACKREST_PORT
		log.Infof("setting BackrestPort to default %d", c.Cluster.BackrestPort)
	}
	if c.Cluster.PGBadgerPort == "" {
		c.Cluster.PGBadgerPort = DEFAULT_PGBADGER_PORT
		log.Infof("setting PGBadgerPort to default %s", c.Cluster.PGBadgerPort)
	} else {
		if _, err := strconv.Atoi(c.Cluster.PGBadgerPort); err != nil {
			return errors.New(errPrefix + "Invalid PGBadgerPort: " + err.Error())
		}
	}
	if c.Cluster.ExporterPort == "" {
		c.Cluster.ExporterPort = DEFAULT_EXPORTER_PORT
		log.Infof("setting ExporterPort to default %s", c.Cluster.ExporterPort)
	} else {
		if _, err := strconv.Atoi(c.Cluster.ExporterPort); err != nil {
			return errors.New(errPrefix + "Invalid ExporterPort: " + err.Error())
		}
	}
	if c.Cluster.Port == "" {
		c.Cluster.Port = DEFAULT_POSTGRES_PORT
		log.Infof("setting Postgres Port to default %s", c.Cluster.Port)
	} else {
		if _, err := strconv.Atoi(c.Cluster.Port); err != nil {
			return errors.New(errPrefix + "Invalid Port: " + err.Error())
		}
	}

	{
		storageNotDefined := func(setting, value string) error {
			return fmt.Errorf("%s%s setting is invalid: %q is not defined", errPrefix, setting, value)
		}
		if _, ok := c.Storage[c.PrimaryStorage]; !ok {
			return storageNotDefined("PrimaryStorage", c.PrimaryStorage)
		}
		if _, ok := c.Storage[c.BackrestStorage]; !ok {
			log.Warning("BackrestStorage setting not set, will use PrimaryStorage setting")
			c.Storage[c.BackrestStorage] = c.Storage[c.PrimaryStorage]
		}
		if _, ok := c.Storage[c.BackupStorage]; !ok {
			return storageNotDefined("BackupStorage", c.BackupStorage)
		}
		if _, ok := c.Storage[c.ReplicaStorage]; !ok {
			return storageNotDefined("ReplicaStorage", c.ReplicaStorage)
		}
		if _, ok := c.Storage[c.WALStorage]; c.WALStorage != "" && !ok {
			return storageNotDefined("WALStorage", c.WALStorage)
		}
		for k := range c.Storage {
			_, err = c.GetStorageSpec(k)
			if err != nil {
				return err
			}
		}
	}

	if c.Pgo.PGOImagePrefix == "" {
		return errors.New(errPrefix + "Pgo.PGOImagePrefix is required")
	}
	if c.Pgo.PGOImageTag == "" {
		return errors.New(errPrefix + "Pgo.PGOImageTag is required")
	}

	if c.Cluster.ServiceType == "" {
		log.Warn("Cluster.ServiceType not set, using default, ClusterIP ")
		c.Cluster.ServiceType = DEFAULT_SERVICE_TYPE
	} else {
		if c.Cluster.ServiceType != DEFAULT_SERVICE_TYPE &&
			c.Cluster.ServiceType != LOAD_BALANCER_SERVICE_TYPE &&
			c.Cluster.ServiceType != NODEPORT_SERVICE_TYPE {
			return errors.New(errPrefix + "Cluster.ServiceType is required to be either ClusterIP, NodePort, or LoadBalancer")
		}
	}

	if c.Cluster.CCPImagePrefix == "" {
		return errors.New(errPrefix + "Cluster.CCPImagePrefix is required")
	}

	if c.Cluster.CCPImageTag == "" {
		return errors.New(errPrefix + "Cluster.CCPImageTag is required")
	}

	if c.Cluster.User == "" {
		return errors.New(errPrefix + "Cluster.User is required")
	} else {
		// validates that username can be used as the kubernetes secret name
		// Must consist of lower case alphanumeric characters,
		// '-' or '.', and must start and end with an alphanumeric character
		errs := validation.IsDNS1123Subdomain(c.Cluster.User)
		if len(errs) > 0 {
			var msg string
			for i := range errs {
				msg = msg + errs[i]
			}
			return errors.New(errPrefix + msg)
		}

		// validate any of the resources and if they are unavailable, set defaults
		if c.Cluster.DefaultInstanceResourceMemory.IsZero() {
			c.Cluster.DefaultInstanceResourceMemory = DefaultInstanceResourceMemory
		}

		log.Infof("default instance memory set to [%s]", c.Cluster.DefaultInstanceResourceMemory.String())

		if c.Cluster.DefaultBackrestResourceMemory.IsZero() {
			c.Cluster.DefaultBackrestResourceMemory = DefaultBackrestResourceMemory
		}

		log.Infof("default pgbackrest repository memory set to [%s]", c.Cluster.DefaultBackrestResourceMemory.String())

		if c.Cluster.DefaultPgBouncerResourceMemory.IsZero() {
			c.Cluster.DefaultPgBouncerResourceMemory = DefaultPgBouncerResourceMemory
		}

		log.Infof("default pgbouncer memory set to [%s]", c.Cluster.DefaultPgBouncerResourceMemory.String())
	}

	// if provided, ensure that the type of pod anti-affinity values are valid
	podAntiAffinityType := crv1.PodAntiAffinityType(c.Cluster.PodAntiAffinity)
	if err := podAntiAffinityType.Validate(); err != nil {
		return errors.New(errPrefix + "Invalid value provided for Cluster.PodAntiAffinityType")
	}

	podAntiAffinityType = crv1.PodAntiAffinityType(c.Cluster.PodAntiAffinityPgBackRest)
	if err := podAntiAffinityType.Validate(); err != nil {
		return errors.New(errPrefix + "Invalid value provided for Cluster.PodAntiAffinityPgBackRest")
	}

	podAntiAffinityType = crv1.PodAntiAffinityType(c.Cluster.PodAntiAffinityPgBouncer)
	if err := podAntiAffinityType.Validate(); err != nil {
		return errors.New(errPrefix + "Invalid value provided for Cluster.PodAntiAffinityPgBouncer")
	}

	return err
}

// GetPodAntiAffinitySpec accepts possible user-defined values for what the
// pod anti-affinity spec should be, which include rules for:
// - PostgreSQL instances
// - pgBackRest
// - pgBouncer
func (c *PgoConfig) GetPodAntiAffinitySpec(cluster, pgBackRest, pgBouncer crv1.PodAntiAffinityType) (crv1.PodAntiAffinitySpec, error) {
	spec := crv1.PodAntiAffinitySpec{}

	// first, set the values for the PostgreSQL cluster, which is the "default"
	// value. Otherwise, set the default to that in the configuration
	if cluster != "" {
		spec.Default = cluster
	} else {
		spec.Default = crv1.PodAntiAffinityType(c.Cluster.PodAntiAffinity)
	}

	// perform a validation check against the default type
	if err := spec.Default.Validate(); err != nil {
		log.Error(err)
		return spec, err
	}

	// now that the default is set, determine if the user or the configuration
	// overrode the settings for pgBackRest and pgBouncer. The heuristic is as
	// such:
	//
	// 1. If the user provides a value, use that value
	// 2. If there is a value provided in the configuration, use that value
	// 3. If there is a value in the cluster default, use that value, which also
	//    encompasses using the default value in the config at this point in the
	//    execution.
	//
	// First, do pgBackRest:
	switch {
	case pgBackRest != "":
		spec.PgBackRest = pgBackRest
	case c.Cluster.PodAntiAffinityPgBackRest != "":
		spec.PgBackRest = crv1.PodAntiAffinityType(c.Cluster.PodAntiAffinityPgBackRest)
	case spec.Default != "":
		spec.PgBackRest = spec.Default
	}

	// perform a validation check against the pgBackRest type
	if err := spec.PgBackRest.Validate(); err != nil {
		log.Error(err)
		return spec, err
	}

	// Now, pgBouncer:
	switch {
	case pgBouncer != "":
		spec.PgBouncer = pgBouncer
	case c.Cluster.PodAntiAffinityPgBackRest != "":
		spec.PgBouncer = crv1.PodAntiAffinityType(c.Cluster.PodAntiAffinityPgBouncer)
	case spec.Default != "":
		spec.PgBouncer = spec.Default
	}

	// perform a validation check against the pgBackRest type
	if err := spec.PgBouncer.Validate(); err != nil {
		log.Error(err)
		return spec, err
	}

	return spec, nil
}

func (c *PgoConfig) GetStorageSpec(name string) (crv1.PgStorageSpec, error) {
	var err error
	storage := crv1.PgStorageSpec{}

	s, ok := c.Storage[name]
	if !ok {
		err = errors.New("invalid Storage name " + name)
		log.Error(err)
		return storage, err
	}

	storage.StorageClass = s.StorageClass
	storage.AccessMode = s.AccessMode
	storage.Size = s.Size
	storage.StorageType = s.StorageType
	storage.MatchLabels = s.MatchLabels
	storage.SupplementalGroups = s.SupplementalGroups

	if storage.MatchLabels != "" {
		test := strings.Split(storage.MatchLabels, "=")
		if len(test) != 2 {
			err = errors.New("invalid Storage config " + name + " MatchLabels needs to be in key=value format.")
			log.Error(err)
			return storage, err
		}
	}

	return storage, err

}

func (c *PgoConfig) GetConfig(clientset kubernetes.Interface, namespace string) error {

	cMap, rootPath := getRootPath(clientset, namespace)

	var yamlFile []byte
	var err error

	//get the pgo.yaml config file
	if cMap != nil {
		str := cMap.Data[CONFIG_PATH]
		if str == "" {
			errMsg := fmt.Sprintf("could not get %s from ConfigMap", CONFIG_PATH)
			return errors.New(errMsg)
		}
		yamlFile = []byte(str)
	} else {
		yamlFile, err = ioutil.ReadFile(rootPath + CONFIG_PATH)
		if err != nil {
			log.Errorf("yamlFile.Get err   #%v ", err)
			return err
		}
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Errorf("Unmarshal: %v", err)
		return err
	}

	// validate the pgo.yaml config file
	if err := c.Validate(); err != nil {
		log.Error(err)
		return err
	}

	c.CheckEnv()

	//load up all the templates
	PgoDefaultServiceAccountTemplate, err = c.LoadTemplate(cMap, rootPath, PGODefaultServiceAccountPath)
	if err != nil {
		return err
	}
	PgoBackrestServiceAccountTemplate, err = c.LoadTemplate(cMap, rootPath, PGOBackrestServiceAccountPath)
	if err != nil {
		return err
	}
	PgoTargetServiceAccountTemplate, err = c.LoadTemplate(cMap, rootPath, PGOTargetServiceAccountPath)
	if err != nil {
		return err
	}
	PgoTargetRoleBindingTemplate, err = c.LoadTemplate(cMap, rootPath, PGOTargetRoleBindingPath)
	if err != nil {
		return err
	}
	PgoBackrestRoleTemplate, err = c.LoadTemplate(cMap, rootPath, PGOBackrestRolePath)
	if err != nil {
		return err
	}
	PgoBackrestRoleBindingTemplate, err = c.LoadTemplate(cMap, rootPath, PGOBackrestRoleBindingPath)
	if err != nil {
		return err
	}
	PgoTargetRoleTemplate, err = c.LoadTemplate(cMap, rootPath, PGOTargetRolePath)
	if err != nil {
		return err
	}
	PgoPgServiceAccountTemplate, err = c.LoadTemplate(cMap, rootPath, PGOPgServiceAccountPath)
	if err != nil {
		return err
	}
	PgoPgRoleTemplate, err = c.LoadTemplate(cMap, rootPath, PGOPgRolePath)
	if err != nil {
		return err
	}
	PgoPgRoleBindingTemplate, err = c.LoadTemplate(cMap, rootPath, PGOPgRoleBindingPath)
	if err != nil {
		return err
	}

	PVCTemplate, err = c.LoadTemplate(cMap, rootPath, pvcPath)
	if err != nil {
		return err
	}

	PolicyJobTemplate, err = c.LoadTemplate(cMap, rootPath, policyJobTemplatePath)
	if err != nil {
		return err
	}

	ContainerResourcesTemplate, err = c.LoadTemplate(cMap, rootPath, containerResourcesTemplatePath)
	if err != nil {
		return err
	}

	PgoBackrestRepoServiceTemplate, err = c.LoadTemplate(cMap, rootPath, pgoBackrestRepoServiceTemplatePath)
	if err != nil {
		return err
	}

	PgoBackrestRepoTemplate, err = c.LoadTemplate(cMap, rootPath, pgoBackrestRepoTemplatePath)
	if err != nil {
		return err
	}

	PgmonitorEnvVarsTemplate, err = c.LoadTemplate(cMap, rootPath, pgmonitorEnvVarsPath)
	if err != nil {
		return err
	}

	PgbackrestEnvVarsTemplate, err = c.LoadTemplate(cMap, rootPath, pgbackrestEnvVarsPath)
	if err != nil {
		return err
	}

	PgbackrestS3EnvVarsTemplate, err = c.LoadTemplate(cMap, rootPath, pgbackrestS3EnvVarsPath)
	if err != nil {
		return err
	}

	PgAdminTemplate, err = c.LoadTemplate(cMap, rootPath, pgAdminTemplatePath)
	if err != nil {
		return err
	}

	PgAdminServiceTemplate, err = c.LoadTemplate(cMap, rootPath, pgAdminServiceTemplatePath)
	if err != nil {
		return err
	}

	PgbouncerTemplate, err = c.LoadTemplate(cMap, rootPath, pgbouncerTemplatePath)
	if err != nil {
		return err
	}

	PgbouncerConfTemplate, err = c.LoadTemplate(cMap, rootPath, pgbouncerConfTemplatePath)
	if err != nil {
		return err
	}

	PgbouncerUsersTemplate, err = c.LoadTemplate(cMap, rootPath, pgbouncerUsersTemplatePath)
	if err != nil {
		return err
	}

	PgbouncerHBATemplate, err = c.LoadTemplate(cMap, rootPath, pgbouncerHBATemplatePath)
	if err != nil {
		return err
	}

	ServiceTemplate, err = c.LoadTemplate(cMap, rootPath, serviceTemplatePath)
	if err != nil {
		return err
	}

	RmdatajobTemplate, err = c.LoadTemplate(cMap, rootPath, rmdatajobPath)
	if err != nil {
		return err
	}

	BackrestjobTemplate, err = c.LoadTemplate(cMap, rootPath, backrestjobPath)
	if err != nil {
		return err
	}

	BackrestRestorejobTemplate, err = c.LoadTemplate(cMap, rootPath, backrestRestorejobPath)
	if err != nil {
		return err
	}

	PgDumpBackupJobTemplate, err = c.LoadTemplate(cMap, rootPath, pgDumpBackupJobPath)
	if err != nil {
		return err
	}

	PgRestoreJobTemplate, err = c.LoadTemplate(cMap, rootPath, pgRestoreJobPath)
	if err != nil {
		return err
	}

	PVCMatchLabelsTemplate, err = c.LoadTemplate(cMap, rootPath, pvcMatchLabelsPath)
	if err != nil {
		return err
	}

	PVCStorageClassTemplate, err = c.LoadTemplate(cMap, rootPath, pvcSCPath)
	if err != nil {
		return err
	}

	AffinityTemplate, err = c.LoadTemplate(cMap, rootPath, affinityTemplatePath)
	if err != nil {
		return err
	}

	PodAntiAffinityTemplate, err = c.LoadTemplate(cMap, rootPath, podAntiAffinityTemplatePath)
	if err != nil {
		return err
	}

	ExporterTemplate, err = c.LoadTemplate(cMap, rootPath, exporterTemplatePath)
	if err != nil {
		return err
	}

	BadgerTemplate, err = c.LoadTemplate(cMap, rootPath, badgerTemplatePath)
	if err != nil {
		return err
	}

	DeploymentTemplate, err = c.LoadTemplate(cMap, rootPath, deploymentTemplatePath)
	if err != nil {
		return err
	}

	BootstrapTemplate, err = c.LoadTemplate(cMap, rootPath, bootstrapTemplatePath)
	if err != nil {
		return err
	}

	return nil
}

func getRootPath(clientset kubernetes.Interface, namespace string) (*v1.ConfigMap, string) {

	cMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), CustomConfigMapName, metav1.GetOptions{})
	if err == nil {
		log.Infof("Config: %s ConfigMap found, using config files from the configmap", CustomConfigMapName)
		return cMap, ""
	}
	log.Infof("Config: %s ConfigMap NOT found, using default baked-in config files from %s", CustomConfigMapName, DefaultConfigsPath)

	return nil, DefaultConfigsPath
}

// LoadTemplate will load a JSON template from a path
func (c *PgoConfig) LoadTemplate(cMap *v1.ConfigMap, rootPath, path string) (*template.Template, error) {
	var value string
	var err error

	// Determine if there exists a configmap entry for the template file.
	if cMap != nil {
		// Get the data that is stored in the configmap
		value = cMap.Data[path]
	}

	// if the configmap does not exist, or there is no data in the configmap for
	// this particular configuration template, attempt to load the template from
	// the default configuration
	if cMap == nil || value == "" {
		value, err = c.DefaultTemplate(path)

		if err != nil {
			return nil, err
		}
	}

	// if we have a value for the templated file, return
	return template.Must(template.New(path).Parse(value)), nil

}

// DefaultTemplate attempts to load a default configuration template file
func (c *PgoConfig) DefaultTemplate(path string) (string, error) {
	// set the lookup value for the file path based on the default configuration
	// path and the template file requested to be loaded
	fullPath := DefaultConfigsPath + path

	log.Debugf("No entry in cmap loading default path [%s]", fullPath)

	// read in the file from the default path
	buf, err := ioutil.ReadFile(fullPath)

	if err != nil {
		log.Errorf("error: could not read %s", fullPath)
		log.Error(err)
		return "", err
	}

	// extract the value of the default configuration file and return
	value := string(buf)

	return value, nil
}

// CheckEnv is mostly used for the OLM deployment use case
// when someone wants to deploy with OLM, use the baked-in
// configuration, but use a different set of images, by
// setting these env vars in the OLM CSV, users can override
// the baked in images
func (c *PgoConfig) CheckEnv() {
	pgoImageTag := os.Getenv("PGO_IMAGE_TAG")
	if pgoImageTag != "" {
		c.Pgo.PGOImageTag = pgoImageTag
		log.Infof("CheckEnv: using PGO_IMAGE_TAG env var: %s", pgoImageTag)
	}
	pgoImagePrefix := os.Getenv("PGO_IMAGE_PREFIX")
	if pgoImagePrefix != "" {
		c.Pgo.PGOImagePrefix = pgoImagePrefix
		log.Infof("CheckEnv: using PGO_IMAGE_PREFIX env var: %s", pgoImagePrefix)
	}
	ccpImageTag := os.Getenv("CCP_IMAGE_TAG")
	if ccpImageTag != "" {
		c.Cluster.CCPImageTag = ccpImageTag
		log.Infof("CheckEnv: using CCP_IMAGE_TAG env var: %s", ccpImageTag)
	}
	ccpImagePrefix := os.Getenv("CCP_IMAGE_PREFIX")
	if ccpImagePrefix != "" {
		c.Cluster.CCPImagePrefix = ccpImagePrefix
		log.Infof("CheckEnv: using CCP_IMAGE_PREFIX env var: %s", ccpImagePrefix)
	}
}
