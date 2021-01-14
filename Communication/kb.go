package Communication

var (
	features string = "api/features"
	space    string = "api/spaces/space"
)

func (comm *Comm) GetFeatures(username, password string) (string, error) {
	return stringifyResponse(comm.Curl(kb, username, password, features, "GET", nil))
}

func (comm *Comm) CreateSpace(username, password string, data []byte) (string, error) {
	return stringifyResponse(comm.Curl(kb, username, password, space, "POST", data))
}

func (comm *Comm) UpdateSpace(username, password string, id string, data []byte) (string, error) {
	path := space + "/" + id
	return stringifyResponse(comm.Curl(kb, username, password, path, "PUT", data))
}

func (comm *Comm) DeleteSpace(username, password, id string) (string, error) {
	path := space + "/" + id
	return stringifyResponse(comm.Curl(kb, username, password, path, "DELETE", nil))
}

func (comm *Comm) GetSpace(username, password, id string) (string, error) {
	path := space + "/" + id
	return stringifyResponse(comm.Curl(kb, username, password, path, "GET", nil))
}

func (comm *Comm) GetAllSpace(username, password string) (string, error) {
	return stringifyResponse(comm.Curl(kb, username, password, space, "GET", nil))
}
