package Communication

import (
	"encoding/json"
	"fmt"
	structs "goclient/RestStruct/Elasticsearch/structs"
)

var (
	createSnapshot    string = "/_snapshot/%s/%s"
	getSnapshot       string = "/_snapshot/%s/%s"
	getSnapshotStatus string = "/_snapshot/%s/%s/_status"
	deleteSnapshot    string = "/_snapshot/%s/%s"
	restoreSnapshot   string = "/_snapshot/%s/%s/_restore"
)

func (comm *Comm) GetConnection(username, password string) (string, error) {
	return stringifyResponse(comm.Curl(es, username, password, "", "GET", nil))
}

func (comm *Comm) CreateSnapshot(
	username, password string,
	request structs.CreateSnapshotRequest,
) (string, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf(createSnapshot, request.My_Repo, request.My_Snapshot)
	return stringifyResponse(comm.Curl(es, username, password, path, "POST", data))
}

func (comm *Comm) GetSnapshot(
	username, password string,
	request structs.GetSnapshotRequest,
) (string, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf(getSnapshot, request.My_Repo, request.My_Snapshot)
	return stringifyResponse(comm.Curl(es, username, password, path, "GET", data))
}

func (comm *Comm) GetSnapshotStatus(
	username, password string,
	request structs.GetSnapshotStatusRequest,
) (string, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf(getSnapshotStatus, request.My_Repo, request.My_Snapshot)
	return stringifyResponse(comm.Curl(es, username, password, path, "GET", data))
}

func (comm *Comm) RestoreSnapshot(
	username, password string,
	request structs.RestoreSnapshotRequest,
) (string, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf(restoreSnapshot, request.My_Repo, request.My_Snapshot)
	return stringifyResponse(comm.Curl(es, username, password, path, "POST", data))
}

func (comm *Comm) DeleteSnapshot(
	username, password string,
	request structs.DeleteSnapshotRequest,
) (string, error) {
	path := fmt.Sprintf(deleteSnapshot, request.My_Repo, request.My_Snapshot)
	return stringifyResponse(comm.Curl(es, username, password, path, "DELETE", nil))
}
