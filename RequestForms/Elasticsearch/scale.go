package Elasticsearch

type ScaleRequest struct {
	Replicas int32 `json:"replicas"`
}
