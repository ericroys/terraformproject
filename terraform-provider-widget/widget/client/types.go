package client

import (
	"fmt"
)

//IsValid interface for validation of object
type IsValid interface {
	IsValid() error
}

//Widget defines a full Widget object returned from api
type Widget struct {
	ID   string `json:"id"`
	UID  string `json:"uid"`
	Name string `json:"name"`
	Size string `json:"size"`
}

//IsValid implements interface to validate Widget
func (w *Widget) IsValid() error {
	if len(w.ID) == 0 ||
		len(w.UID) == 0 ||
		len(w.Name) == 0 {
		return fmt.Errorf("Unable to tranform to Widget due to missing parameters")
	}
	return nil
}

//WidgetNew is partial Widget including only fields
//needed for api create
type WidgetNew struct {
	Name string `json:"name"`
	Size string `json:"size"`
}

//IsValid implements interface to validate WidgetNew
func (w *WidgetNew) IsValid() error {
	if len(w.Name) == 0 {
		return fmt.Errorf("Missing Name")
	}
	return nil
}
