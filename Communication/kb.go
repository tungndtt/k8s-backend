package Communication

var (
	features string = "api/features"
	space    string = "api/spaces/space"
)

func (comm *Comm) GetFeatures() (string, error) {
	return stringifyResponse(comm.Curl(kb, features, "GET", nil))
}

func (comm *Comm) CreateSpace(data []byte) (string, error) {
	return stringifyResponse(comm.Curl(kb, space, "POST", data))
}

func (comm *Comm) UpdateSpace(id string, data []byte) (string, error) {
	path := space + "/" + id
	return stringifyResponse(comm.Curl(kb, path, "PUT", data))
}

func (comm *Comm) DeleteSpace(id string) (string, error) {
	path := space + "/" + id
	return stringifyResponse(comm.Curl(kb, path, "DELETE", nil))
}

func (comm *Comm) GetSpace(id string) (string, error) {
	path := space + "/" + id
	return stringifyResponse(comm.Curl(kb, path, "GET", nil))
}

func (comm *Comm) GetAllSpace() (string, error) {
	return stringifyResponse(comm.Curl(kb, space, "GET", nil))
}
