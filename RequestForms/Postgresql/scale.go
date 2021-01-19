package Postgresql

import (
	structs "goclient/RestStruct/Postgres/structs"
)

type ScaleRequest struct {
	Name string `json:"name"`
	structs.ClusterScaleRequest
}
