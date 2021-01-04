package Communication

func (comm *Comm) GetConnection(username, password string, port int32) (string, error) {
	return stringifyResponse(comm.Curl(es, username, password, "", "GET", port, nil))
}
