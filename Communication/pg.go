package Communication

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	structs "goclient/RestStruct/Postgres/structs"
)

var (
	version         string = "version"
	showclusters    string = "showclusters"
	createclusters  string = "clusters"
	updateclusters  string = "clustersupdate"
	deleteclusters  string = "clustersdelete"
	createbackrest  string = "backrestbackup"
	backrest        string = "backrest"
	showBackup      string = "%s/%s?version=%s&namespace=%s&selector=%s"
	createpgdump    string = "pgdumpbackup"
	pgdump          string = "pgdump"
	restorebackrest string = "restore"
	restorepgdump   string = "pgdumprestore"
	scaleCluster    string = "clusters/scale/%s"
	scaleDown       string = "scaledown/%s?version=%s&namespace=%s&replica-name=%s&delete-data=%s"
	scale           string = "scale/%s?version=%s&namespace=%s"
)

// parse the response body into the given form
func parseInto(resp *http.Response, form interface{}) error {
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return json.NewDecoder(resp.Body).Decode(form)
	} else {
		return errors.Unwrap(fmt.Errorf("Status: %s, StatusCode: %d", resp.Status, resp.StatusCode))
	}
}

func (comm *Comm) GetVersion() (*structs.VersionResponse, error) {
	resp, err := comm.Curl(pg, version, "GET", nil)
	if err != nil {
		return nil, err
	}
	var res structs.VersionResponse
	err = parseInto(resp, &res)
	return &res, err
}

// clusters

// show clusters based on request
func (comm *Comm) ShowClusters(request *structs.ShowClusterRequest) (*structs.ShowClusterResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, showclusters, "POST", body)
	if err != nil {
		return nil, err
	}
	var res structs.ShowClusterResponse
	err = parseInto(resp, &res)
	return &res, err
}

// create cluster based on request
func (comm *Comm) CreateClusters(request *structs.CreateClusterRequest) (*structs.CreateClusterResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, createclusters, "POST", body)
	if err != nil {
		return nil, err
	}
	var res structs.CreateClusterResponse
	err = parseInto(resp, &res)
	return &res, err
}

// update cluster based on request
func (comm *Comm) UpdateClusters(request *structs.UpdateClusterRequest) (*structs.UpdateClusterResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, updateclusters, "POST", body)
	if err != nil {
		return nil, err
	}
	var res structs.UpdateClusterResponse
	err = parseInto(resp, &res)
	return &res, err
}

// delete cluster based on request
func (comm *Comm) DeleteClusters(request *structs.DeleteClusterRequest) (*structs.DeleteClusterResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, deleteclusters, "POST", body)
	if err != nil {
		return nil, err
	}
	var res structs.DeleteClusterResponse
	err = parseInto(resp, &res)
	return &res, err
}

// backrest backup

// create backrest backup based on request
func (comm *Comm) CreateBackrest(request *structs.CreateBackrestBackupRequest) (*structs.CreateBackrestBackupResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, createbackrest, "POST", body)
	if err != nil {
		return nil, err
	}
	var res structs.CreateBackrestBackupResponse
	err = parseInto(resp, &res)
	return &res, err
}

// delete backrest backup based on request
func (comm *Comm) DeleteBackrest(request *structs.DeleteBackrestBackupRequest) (*structs.DeleteBackrestBackupResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, backrest, "DELETE", body)
	if err != nil {
		return nil, err
	}
	var res structs.DeleteBackrestBackupResponse
	err = parseInto(resp, &res)
	return &res, err
}

// show backrest backup based on request
func (comm *Comm) ShowBackrest(namespace, cluster_name, version, selector string) (*structs.ShowBackrestResponse, error) {
	path := fmt.Sprintf(showBackup, backrest, cluster_name, version, namespace, selector)
	resp, err := comm.Curl(pg, path, "GET", nil)
	if err != nil {
		return nil, err
	}
	var res structs.ShowBackrestResponse
	err = parseInto(resp, &res)
	return &res, err
}

// restore backrest backup based on request
func (comm *Comm) RestoreBackrest(request *structs.RestoreRequest) (*structs.RestoreResponse, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, restorebackrest, "POST", data)
	if err != nil {
		return nil, err
	}
	var res structs.RestoreResponse
	err = parseInto(resp, &res)
	return &res, err
}

// pgdump backup

// create dump backup based on request
func (comm *Comm) CreatePgDumpBackup(request *structs.CreatepgDumpBackupRequest) (*structs.CreatepgDumpBackupResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, createpgdump, "POST", body)
	if err != nil {
		return nil, err
	}
	var res structs.CreatepgDumpBackupResponse
	err = parseInto(resp, &res)
	return &res, err
}

// show dump backup based on request
func (comm *Comm) ShowPgDumpBackup(namespace, cluster_name, version, selector string) (*structs.ShowBackupResponse, error) {
	path := fmt.Sprintf(showBackup, pgdump, cluster_name, version, namespace, selector)
	resp, err := comm.Curl(pg, path, "GET", nil)
	if err != nil {
		return nil, err
	}
	var res structs.ShowBackupResponse
	err = parseInto(resp, &res)
	return &res, err
}

// restore backrest backup based on request
func (comm *Comm) RestorePgDump(request *structs.PgRestoreRequest) (*structs.PgRestoreResponse, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, restorepgdump, "POST", data)
	if err != nil {
		return nil, err
	}
	var res structs.PgRestoreResponse
	err = parseInto(resp, &res)
	return &res, err
}

// scale
func (comm *Comm) ScaleCluster(name string, request *structs.ClusterScaleRequest) (*structs.ClusterScaleResponse, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, fmt.Sprintf(scaleCluster, name), "POST", data)
	if err != nil {
		return nil, err
	}
	var res structs.ClusterScaleResponse
	err = parseInto(resp, &res)
	return &res, err
}

func (comm *Comm) ScaleDown(name, version, namespace, replica_name, delete_data string) (*structs.ScaleDownResponse, error) {
	resp, err := comm.Curl(
		pg,
		fmt.Sprintf(
			scaleDown,
			name, version, namespace, replica_name, delete_data,
		),
		"GET", nil,
	)
	if err != nil {
		return nil, err
	}
	var res structs.ScaleDownResponse
	err = parseInto(resp, &res)
	return &res, err
}

func (comm *Comm) GetScaleQuery(name, version, namespace string) (*structs.ScaleQueryResponse, error) {
	resp, err := comm.Curl(
		pg,
		fmt.Sprintf(
			scale,
			name,
			version,
			namespace,
		),
		"GET", nil,
	)
	if err != nil {
		return nil, err
	}
	var res structs.ScaleQueryResponse
	err = parseInto(resp, &res)
	return &res, err
}
