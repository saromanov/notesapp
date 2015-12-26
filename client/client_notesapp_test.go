package client

import(
  "testing"
  "fmt"

)

func TestCreateNote(t *testing.T) {
	cli := ClientNotesapp{Addr:"testing"}
	err := cli.CreateNote("first", "second")
	if err != nil {
		t.Errorf(fmt.Sprintf("%v", err))
	}
}

func TestGetAllNotes(t *testing.T) {
	cli := ClientNotesapp{Addr:"testing"}
	_, err := cli.GetAllNotes()
	if err != nil {
		t.Errorf(fmt.Sprintf("%v", err))
	}
}

func TestRemoveNote(t *testing.T) {
	cli := ClientNotesapp{Addr:"testing"}
	err := cli.RemoveNote("test")
	if err != nil {
		t.Errorf(fmt.Sprintf("%v", err))
	}
}

func TestUpdateNote(t *testing.T) {
	cli := ClientNotesapp{Addr:"testing"}
	err := cli.UpdateNote("test1", "test2", "new note")
	if err != nil {
		t.Errorf(fmt.Sprintf("%v", err))
	}
}