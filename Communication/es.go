package Communication

func (comm *Comm) GetConnection(username, password string) (string, error) {
	return stringifyResponse(comm.Curl(es, username, password, "", "GET", nil))
}
