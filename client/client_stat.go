package client


// IncGets provides increment on stat microservice
func (cli *ClientNotesapp) IncGets() error {
	var err error
	note := InsertNoteRequest{Title: "test"}

	_, err = request(cli.Addr, "POST", note)
	if err != nil {
		return err
	}

	return nil
}