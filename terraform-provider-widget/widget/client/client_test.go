package client

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalWidget(t *testing.T) {

	in := `{"id": "123", "name": "testWidget", "size": "enourmous", "uid": "f45fdkksj89g"}`
	inBad := `{"id": "123", "badName": "testWidget", "size": "enourmous", "uid": "f45fdkksj89g"}`
	w := Widget{}

	err := json.Unmarshal([]byte(in), &w)

	if err != nil {
		t.Fatal("Failed to unmarshal widget", err)
	}
	t.Logf("%v", w)

	err = json.Unmarshal([]byte(inBad), &w)

	if err == nil {
		t.Fatal("Expected error but got none")
	}
	t.Log(err)
}

func TestCreateWidget(t *testing.T) {
	w := WidgetNew{
		Name: "test_time",
		Size: "not very",
	}

	c, err := NewClient("http://localhost:8080/api")
	if err != nil {
		t.Fatal(err)
	}
	x, err := c.CreateWidget(w)

	if err != nil {
		t.Fatal("ERROR: ", err.Error())
	}

	err = x.IsValid()
	if err != nil {
		t.Fatal("ERROR: ", err.Error())
	}
}

func TestDeleteWidget(t *testing.T) {
	c, _ := NewClient("http://localhost:8080/api")
	ids := []string{
		"1", "2", "3", "4", "unknown",
	}

	for i, id := range ids {
		err := c.DeleteWidget(id)
		if i < 4 {
			if err != nil {
				t.Fatalf("Not able to delete id: %s due to error: %s", id, err)
			}
		} else {
			if err == nil {
				t.Fatal("Expected error Not Found/Not Deleted, Got none")
			}
		}
	}
}

func TestUpdateWidget(t *testing.T) {
	c, _ := NewClient("http://localhost:8080/api")

	w := WidgetNew{
		Name: "newWidget1",
		Size: "newSize",
	}
	i := `1`

	x, err := c.UpdateWidget(i, w)
	if err != nil {
		t.Fatal("ERROR: ", err.Error())
	}
	t.Logf("Update returned: %v", x)
}
func TestGetWidget(t *testing.T) {

	ids := []string{
		"1", "2", "3", "4", "unknown",
	}

	c, _ := NewClient("http://localhost:8080/api")

	for i, id := range ids {
		x, err := c.GetWidget(id)

		//only "unknown should fail"
		if i < 4 {
			if err != nil {
				t.Fatal("ERROR: ", err.Error())
			}
			if x.ID != id {
				t.Fatalf("Expected id: %s, retrieved: %s", id, x.ID)
			}
		} else {
			if err == nil {
				t.Fatal("Expected error Not Found, Got none")
			}
		}

	}

	//t.Fatalf("RESPONSE: %v", out)

}
