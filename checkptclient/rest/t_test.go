package rest

import (
	"testing"
)

// func TestBuilder(t *testing.T) {

// 	b := NewBuilder("http://localhost").
// 		Method(GET).
// 		ContentType("application/json").
// 		ErrorHandler(DefaultErrorHandler{}).Build()
// }

func TestF(t *testing.T) {
	var m map[string]string

	t.Log(m)
	m = make(map[string]string)
	m["x"] = "0"
	t.Log(m)
}
