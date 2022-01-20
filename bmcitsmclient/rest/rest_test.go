package rest

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func getClient() *http.Client {
	trans := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}
	return &http.Client{
		Transport: trans,
		Timeout:   2 * time.Second,
	}
}
func TestRequestable(t *testing.T) {

	r, err := NewRequestBuilder("http://localhost:8080/api/widget/ID:621f812b-64a3-4c13-8f4d-295d3371a91d", getClient()).
		//Auth(AuthBasic{user: "bob", pass: "xxxxx"}).
		Auth(AuthNoAuth{}).
		ContentType("application/json").
		Method(GET).Build()

	if err != nil {
		t.Fatal(err)
	}
	x, err := r.Send()
	fmt.Printf("Content: %s\nError: %v", x, err)
}

func TestRequestablePost(t *testing.T) {
	var msg = (`{"name": "wigglyWidget", "size": "all over the place"}`)
	var m = []byte(msg)

	r, err := NewRequestBuilder("http://localhost:8080/api/widget/", getClient()).
		Auth(AuthNoAuth{}).
		ContentType("application/json").
		Message(m).
		Method(POST).Build()

	if err != nil {
		t.Fatal(err)
	}
	x, err := r.Send()
	fmt.Printf("Content: %s\nError: %v", x, err)
}

func TestRequestableLogin(t *testing.T) {
	var msg = (`{"user": "admin", "password": "vpn12345"}`)
	var m = []byte(msg)

	r, err := NewRequestBuilder("https://10.10.13.96/web_api/v1.4/login", getClient()).
		Auth(AuthNoAuth{}).
		ContentType("application/json").
		Message(m).
		Method(POST).Build()

	if err != nil {
		t.Fatal(err)
	}
	x, err := r.Send()
	fmt.Printf("Content: %s\nError: %v", x, err)
}

// func TestRequest(t *testing.T) {
// 	r := Request{
// 		Headers: map[string]string{
// 			"h1": "h1V",
// 			"h2": "h2v",
// 		},

// 		ContentType: "application/json",
// 		Method:      POST,
// 		URL:         "http://localhost/something",
// 		Auth:        AuthNoAuth{},
// 	}
// 	err := r.Build()

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log("Validation Successful")

// 	t.Logf("Method %s\n", r.Method)
// 	t.Logf("Headers: %v", r.Headers)

// 	rg := RequestGet{
// 		R: Request{
// 			Headers: map[string]string{
// 				"h1": "h1V",
// 				"h2": "h2v",
// 			},

// 			ContentType: "application/json",
// 			Method:      POST,
// 			URL:         "http://localhost/something",
// 			Auth:        AuthNoAuth{},
// 		},
// 	}
// }
