package Postgresql

import (
	structs "goclient/RestStruct/Postgres/structs"
)

type CreateBackupRequest struct {
	structs.CreateBackrestBackupRequest
}

type ShowBackupRequest struct {
	Namespace   string `json:"namespace"`
	Clustername string `json:"cluster_name"`
	Version     string `json:"version"`
	Selector    string `json:"selector"`
}

type DeleteBackupRequest struct {
	structs.DeleteBackrestBackupRequest
}
