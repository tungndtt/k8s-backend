package Elasticsearch

type CreateRepoDto struct {
	Master_timeout string      `json:"master_timeout"`
	Timeout        string      `json:"timeout"`
	Type           string      `json:"type"`
	Settings       RepoSetting `json:"settings"`
	Verify         bool        `json:"verify"`
}

type RepoSetting struct {
	Chunk_size                 string `json:"chunk_size"`
	Compress                   bool   `json:"compress"`
	Max_number_of_snapshots    int16  `json:"max_number_of_snapshots"`
	Max_restore_bytes_per_sec  string `json:"max_restore_bytes_per_sec"`
	Max_snapshot_bytes_per_sec string `json:"max_snapshot_bytes_per_sec"`
	Readonly                   bool   `json:"readonly"`
	Location                   string `json:"location"`
	Delegate_type              string `json:"delegate_type"`
	Url                        string `json:"url"`
}

type RepoDto struct {
	Repository string
	Body       CreateRepoDto
}
