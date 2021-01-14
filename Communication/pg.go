package Communication

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	msgs "goclient/apiservermsgs"
)

var (
	version         string = "version"
	showclusters    string = "showclusters"
	createclusters  string = "clusters"
	updateclusters  string = "clustersupdate"
	deleteclusters  string = "clustersdelete"
	createbackrest  string = "backrestbackup"
	backrest        string = "backrest"
	createpgdump    string = "pgdumpbackup"
	pgdump          string = "pgdump"
	restorebackrest string = "restore"
	restorepgdump   string = "pgdumprestore"
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

func (comm *Comm) GetVersion(username, password string) (*msgs.VersionResponse, error) {
	resp, err := comm.Curl(pg, username, password, version, "GET", nil)
	if err != nil {
		return nil, err
	}
	var res msgs.VersionResponse
	err = parseInto(resp, &res)
	return &res, err
}

// clusters

// show clusters based on request
func (comm *Comm) ShowClusters(username, password string, request msgs.ShowClusterRequest) (*msgs.ShowClusterResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, username, password, showclusters, "POST", body)
	if err != nil {
		return nil, err
	}
	var res msgs.ShowClusterResponse
	err = parseInto(resp, &res)
	return &res, err
}

// create cluster based on request
func (comm *Comm) CreateClusters(username, password string, request msgs.CreateClusterRequest) (*msgs.CreateClusterResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, username, password, createclusters, "POST", body)
	if err != nil {
		return nil, err
	}
	var res msgs.CreateClusterResponse
	err = parseInto(resp, &res)
	return &res, err
}

// update cluster based on request
func (comm *Comm) UpdateClusters(username, password string, request msgs.UpdateClusterRequest) (*msgs.UpdateClusterResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, username, password, updateclusters, "POST", body)
	if err != nil {
		return nil, err
	}
	var res msgs.UpdateClusterResponse
	err = parseInto(resp, &res)
	return &res, err
}

// delete cluster based on request
func (comm *Comm) DeleteClusters(username, password string, request msgs.DeleteClusterRequest) (*msgs.DeleteClusterResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, username, password, deleteclusters, "POST", body)
	if err != nil {
		return nil, err
	}
	var res msgs.DeleteClusterResponse
	err = parseInto(resp, &res)
	return &res, err
}

// backrest backup

// create backrest backup based on request
func (comm *Comm) CreateBackrest(username, password string, request msgs.CreateBackrestBackupRequest) (*msgs.CreateBackrestBackupResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, username, password, createbackrest, "POST", body)
	if err != nil {
		return nil, err
	}
	var res msgs.CreateBackrestBackupResponse
	err = parseInto(resp, &res)
	return &res, err
}

// delete backrest backup based on request
func (comm *Comm) DeleteBackrest(username, password string, request msgs.DeleteBackrestBackupRequest) (*msgs.DeleteBackrestBackupResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, username, password, backrest, "DELETE", body)
	if err != nil {
		return nil, err
	}
	var res msgs.DeleteBackrestBackupResponse
	err = parseInto(resp, &res)
	return &res, err
}

// show backrest backup based on request
func (comm *Comm) ShowBackrest(username, password string, namespace, cluster_name, version, selector string) (*msgs.ShowBackrestResponse, error) {
	path := fmt.Sprintf("%s/%s?version=%s&namespace=%s&selector=%s", backrest, cluster_name, version, namespace, selector)
	resp, err := comm.Curl(pg, username, password, path, "GET", nil)
	if err != nil {
		return nil, err
	}
	var res msgs.ShowBackrestResponse
	err = parseInto(resp, &res)
	return &res, err
}

// restore backrest backup based on request
func (comm *Comm) RestoreBackrest(username, password string, request msgs.RestoreRequest) (*msgs.RestoreResponse, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, username, password, restorebackrest, "POST", data)
	if err != nil {
		return nil, err
	}
	var res msgs.RestoreResponse
	err = parseInto(resp, &res)
	return &res, err
}

// pgdump backup

// create dump backup based on request
func (comm *Comm) CreatePgDumpBackup(username, password string, request msgs.CreatepgDumpBackupRequest) (*msgs.CreatepgDumpBackupResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, username, password, createpgdump, "POST", body)
	if err != nil {
		return nil, err
	}
	var res msgs.CreatepgDumpBackupResponse
	err = parseInto(resp, &res)
	return &res, err
}

// show dump backup based on request
func (comm *Comm) ShowPgDumpBackup(username, password string, namespace, cluster_name, version, selector string) (*msgs.ShowBackupResponse, error) {
	path := fmt.Sprintf("%s/%s?version=%s&namespace=%s&selector=%s", pgdump, cluster_name, version, namespace, selector)
	resp, err := comm.Curl(pg, username, password, path, "GET", nil)
	if err != nil {
		return nil, err
	}
	var res msgs.ShowBackupResponse
	err = parseInto(resp, &res)
	return &res, err
}

// restore backrest backup based on request
func (comm *Comm) RestorePgDump(username, password string, request msgs.PgRestoreRequest) (*msgs.PgRestoreResponse, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, username, password, restorepgdump, "POST", data)
	if err != nil {
		return nil, err
	}
	var res msgs.PgRestoreResponse
	err = parseInto(resp, &res)
	return &res, err
}
