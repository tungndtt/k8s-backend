package Communication

var (
	features string = "api/features"
	space    string = "api/spaces/space"
)

func (comm *Comm) GetFeatures(username, password string, port int32) (string, error) {
	return stringifyResponse(comm.Curl(kb, username, password, features, "GET", port, nil))
}

func (comm *Comm) CreateSpace(username, password string, port int32, data []byte) (string, error) {
	return stringifyResponse(comm.Curl(kb, username, password, space, "POST", port, data))
}

func (comm *Comm) UpdateSpace(username, password string, port int32, id string, data []byte) (string, error) {
	path := space + "/" + id
	return stringifyResponse(comm.Curl(kb, username, password, path, "PUT", port, data))
}

func (comm *Comm) DeleteSpace(username, password, id string, port int32) (string, error) {
	path := space + "/" + id
	return stringifyResponse(comm.Curl(kb, username, password, path, "DELETE", port, nil))
}

func (comm *Comm) GetSpace(username, password, id string, port int32) (string, error) {
	path := space + "/" + id
	return stringifyResponse(comm.Curl(kb, username, password, path, "GET", port, nil))
}

func (comm *Comm) GetAllSpace(username, password string, port int32) (string, error) {
	return stringifyResponse(comm.Curl(kb, username, password, space, "GET", port, nil))
}
