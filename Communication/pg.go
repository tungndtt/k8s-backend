package Communication

import (
	"encoding/json"
	"net/http"

	msgs "goclient/apiservermsgs"
)

var (
	version        string = "version"
	showclusters   string = "showclusters"
	createclusters string = "clusters"
	updateclusters string = "clustersupdate"
	deleteclusters string = "clustersdelete"
)

func (comm *Comm) GetVersion(username, password string, port int32) (*msgs.VersionResponse, error) {
	resp, err := comm.Curl(pg, username, password, version, "GET", port, nil)
	if resp.StatusCode == http.StatusOK {
		var res msgs.VersionResponse
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			return nil, err
		}
		return &res, nil
	}
	return nil, nil
}

func (comm *Comm) ShowClusters(username, password string, port int32, request msgs.ShowClusterRequest) (*msgs.ShowClusterResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, username, password, showclusters, "POST", port, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var res msgs.ShowClusterResponse
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			return nil, err
		}
		return &res, nil
	} else {
		return nil, nil
	}
}

func (comm *Comm) CreateClusters(username, password string, port int32, request msgs.CreateClusterRequest) (*msgs.CreateClusterResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, username, password, createclusters, "POST", port, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var res msgs.CreateClusterResponse
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			return nil, err
		}
		return &res, nil
	}
	return nil, nil
}

func (comm *Comm) UpdateClusters(username, password string, port int32, request msgs.UpdateClusterRequest) (*msgs.UpdateClusterResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, username, password, updateclusters, "POST", port, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var res msgs.UpdateClusterResponse
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			return nil, err
		}
		return &res, nil
	}
	return nil, nil
}

func (comm *Comm) DeleteClusters(username, password string, port int32, request msgs.DeleteClusterRequest) (*msgs.DeleteClusterResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := comm.Curl(pg, username, password, deleteclusters, "POST", port, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var res msgs.DeleteClusterResponse
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			return nil, err
		}
		return &res, nil
	}
	return nil, nil
}
