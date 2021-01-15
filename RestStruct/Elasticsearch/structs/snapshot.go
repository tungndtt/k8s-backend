package structs

type CreateSnapshotRequest struct {
	My_Repo              string `json:"my_repositoriy"`
	My_Snapshot          string `json:"my_snapshot"`
	Ignore_Unavailable   bool   `json:"ignore_unavailable"`
	Indices              string `json:"indices"`
	Include_Global_State bool   `json:"include_global_state"`
	Master_Timeout       string `json:"master_timeout"`
	Metadata             string `json:"metadata"`
	Partial              bool   `json:"partial"`
	Wait_For_Completion  bool   `json:"wait_for_completion"`
}

type GetSnapshotRequest struct {
	My_Repo            string `json:"my_repositoriy"`
	My_Snapshot        string `json:"my_snapshot"`
	Ignore_Unavailable bool   `json:"ignore_unavailable"`
	Verbose            bool   `json:"verbose"`
}

type GetSnapshotStatusRequest struct {
	My_Repo            string `json:"my_repositoriy"`
	My_Snapshot        string `json:"my_snapshot"`
	Ignore_Unavailable bool   `json:"ignore_unavailable"`
}

type RestoreSnapshotRequest struct {
	My_Repo               string `json:"my_repositoriy"`
	My_Snapshot           string `json:"my_snapshot"`
	Ignore_Unavailable    bool   `json:"ignore_unavailable"`
	Ignore_Index_Settings string `json:"ignore_index_settings"`
	Include_Aliases       bool   `json:"include_aliases"`
	Include_Global_State  bool   `json:"include_global_state"`
	Index_Settings        string `json:"index_settings"`
	Indices               string `json:"indices"`
	Partial               bool   `json:"partial"`
	Rename_Pattern        string `json:"rename_pattern"`
	Rename_Replacement    string `json:"rename_replacement"`
	Wait_For_Completion   bool   `json:"wait_for_completion"`
}

type DeleteSnapshotRequest struct {
	My_Repo     string `json:"my_repositoriy"`
	My_Snapshot string `json:"my_snapshot"`
}
