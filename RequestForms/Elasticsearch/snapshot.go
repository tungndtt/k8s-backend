package Elasticsearch

import "goclient/RestStruct/Elasticsearch/structs"

type CreateSnapshotRequest struct {
	My_Repo     string `json:"my_repositoriy"`
	My_Snapshot string `json:"my_snapshot"`
	structs.CreateSnapshot
}

type GetSnapshotRequest struct {
	My_Repo     string `json:"my_repositoriy"`
	My_Snapshot string `json:"my_snapshot"`
	structs.GetSnapshot
}

type GetSnapshotStatusRequest struct {
	My_Repo     string `json:"my_repositoriy"`
	My_Snapshot string `json:"my_snapshot"`
	structs.GetSnapshotStatus
}

type RestoreSnapshotRequest struct {
	My_Repo     string `json:"my_repositoriy"`
	My_Snapshot string `json:"my_snapshot"`
	structs.RestoreSnapshot
}

type DeleteSnapshotRequest struct {
	My_Repo     string `json:"my_repositoriy"`
	My_Snapshot string `json:"my_snapshot"`
}
