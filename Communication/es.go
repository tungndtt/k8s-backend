package Communication

import (
	"encoding/json"
	"fmt"
	forms "goclient/RequestForms/Elasticsearch"
)

var (
	repository        string = "_snapshot/%s"
	createSnapshot    string = "_snapshot/%s/%s"
	getSnapshot       string = "_snapshot/%s/%s"
	getSnapshotStatus string = "_snapshot/%s/%s/_status"
	deleteSnapshot    string = "_snapshot/%s/%s"
	restoreSnapshot   string = "_snapshot/%s/%s/_restore"
)

func (comm *Comm) GetConnection() (string, error) {
	return stringifyResponse(comm.Curl(es, "", "GET", nil))
}

func (comm *Comm) CreateRepository(request forms.RepoDto) (string, error) {
	data, err := json.Marshal(request.Body)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf(repository, request.Repository)
	return stringifyResponse(comm.Curl(es, path, "POST", data))
}

func (comm *Comm) GetRepository(request forms.RepoDto) (string, error) {
	path := fmt.Sprintf(repository, request.Repository)
	return stringifyResponse(comm.Curl(es, path, "GET", nil))
}

func (comm *Comm) CreateSnapshot(request forms.CreateSnapshotRequest) (string, error) {
	data, err := json.Marshal(request.CreateSnapshot)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf(createSnapshot, request.My_Repo, request.My_Snapshot)
	return stringifyResponse(comm.Curl(es, path, "POST", data))
}

func (comm *Comm) GetSnapshot(request forms.GetSnapshotRequest) (string, error) {
	data, err := json.Marshal(request.GetSnapshot)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf(getSnapshot, request.My_Repo, request.My_Snapshot)
	return stringifyResponse(comm.Curl(es, path, "GET", data))
}

func (comm *Comm) GetSnapshotStatus(request forms.GetSnapshotStatusRequest) (string, error) {
	data, err := json.Marshal(request.GetSnapshotStatus)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf(getSnapshotStatus, request.My_Repo, request.My_Snapshot)
	return stringifyResponse(comm.Curl(es, path, "GET", data))
}

func (comm *Comm) RestoreSnapshot(request forms.RestoreSnapshotRequest) (string, error) {
	data, err := json.Marshal(request.RestoreSnapshot)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf(restoreSnapshot, request.My_Repo, request.My_Snapshot)
	return stringifyResponse(comm.Curl(es, path, "POST", data))
}

func (comm *Comm) DeleteSnapshot(request forms.DeleteSnapshotRequest) (string, error) {
	path := fmt.Sprintf(deleteSnapshot, request.My_Repo, request.My_Snapshot)
	return stringifyResponse(comm.Curl(es, path, "DELETE", nil))
}
