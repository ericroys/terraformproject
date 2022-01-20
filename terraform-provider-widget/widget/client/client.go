package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

//Client a client for the Widget API
type Client struct {
	svcurl     string
	httpClient *http.Client
}

//NewClient returns a new initialized Client stucture
//for the Widget API
func NewClient(baseurl string) (*Client, error) {
	ie := `http://localhost:8080/api...`
	if len(baseurl) == 0 {
		return nil, fmt.Errorf("Client expects a base url (i.e. %s)", ie)
	}
	//customize client so set timeout and connection pool stuff
	trans := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}

	return &Client{
		svcurl: baseurl,
		httpClient: &http.Client{
			Transport: trans,
			Timeout:   2 * time.Second,
		},
	}, nil
}

//UpdateWidget updates a Widget by id, with values received in WidgetNew.
//If the widget is not found, and error will indicate so
func (c *Client) UpdateWidget(id string, w WidgetNew) (Widget, error) {
	fmt.Println("WidgetUpdate - Start")
	//Define the Widget as the transformer to use for
	//the unmarshall and what we return to client
	var r Widget
	err := c.update("widget", id, w, &r)

	if err != nil {
		return r, err
	}
	return r, r.IsValid()
}

//CreateWidget creates a Widget object on the service. Expects a WidgetNew
//struct populated for the message and returns a populated Widget. Error is nil
//on successful response.
func (c *Client) CreateWidget(w WidgetNew) (Widget, error) {

	fmt.Println("WidgetCreate -Start")
	//Define the Widget as the transformer to use for
	//the unmarshall and what we return to client
	var r Widget
	err := c.create("widget", w, &r)

	if err != nil {
		return r, err
	}

	return r, r.IsValid()
}

//GetWidget performs a search against the service for a particular
//Widget id and if found, returns a full Widget struct. If nothing is
//found an error is returned
func (c *Client) GetWidget(id string) (Widget, error) {
	fmt.Println("WidgetGet -Starts")
	var r Widget
	err := c.get("widget", id, &r)
	if err != nil {
		return r, err
	}

	return r, r.IsValid()
}

//DeleteWidget performs a search for a particular Widget id and if found,
//the service will delete the Widget. If not found or not deleted, error
// will not be nil
func (c *Client) DeleteWidget(id string) error {
	err := c.delete("widget", id)
	return err
}

func (c *Client) update(uri, id string, msg, trans interface{}) error {
	//build message
	toMsg, err := getMessage(msg)
	if err != nil {
		return err
	}

	//build url
	url, err := c.getPath(uri, id)
	if err != nil {
		return err
	}

	//setup the request
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(toMsg))
	request.Header.Add("Content-Type", "application/json")

	//send the request
	data, err := c.send(request)
	if err != nil {
		return err
	}
	//convert return message
	err = json.Unmarshal([]byte(data), &trans)
	log.Printf("Unmarshal Trans [%v], Error [%v]", trans, err)
	if err != nil {
		return err
	}
	return nil
}

//delete is generalized rest delete function that takes in a uri and an object id
//if successful returns nil error, otherwise error will contain message of failure
func (c *Client) delete(uri string, id string) error {
	//build url
	url, err := c.getPath(uri, id)
	if err != nil {
		return err
	}

	//setup the request
	request, _ := http.NewRequest("DELETE", url, nil)
	request.Header.Add("Content-Type", "application/json")
	_, er := c.send(request)

	return er
}

//get is generalized rest get function that takes in uri, object id, and transformation
//interface for the response.
func (c *Client) get(uri string, id string, trans interface{}) error {
	//build url
	url, err := c.getPath(uri, id)
	if err != nil {
		return err
	}

	//setup the request
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Content-Type", "application/json")

	//send the request
	data, err := c.send(request)
	if err != nil {
		return err
	}
	//convert return message
	err = json.Unmarshal([]byte(data), &trans)
	log.Printf("Unmarshal Trans [%v], Error [%v]", trans, err)
	if err != nil {
		return err
	}
	return nil

}

//create is generalized rest create function that takes in the uri for the service,
//a message interface (must be marshallable to json), and an interface for transforming
//the response body to applicable format.
func (c *Client) create(uri string, msg interface{}, trans interface{}) error {
	//build message
	toMsg, err := getMessage(msg)
	if err != nil {
		return err
	}

	//build url
	url, err := c.getPath(uri, "")
	if err != nil {
		return err
	}

	//setup the request
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(toMsg))
	request.Header.Add("Content-Type", "application/json")

	//send the request
	data, err := c.send(request)
	if err != nil {
		return err
	}
	//convert return message
	err = json.Unmarshal([]byte(data), &trans)
	log.Printf("Unmarshal Trans [%v], Error [%v]", trans, err)
	if err != nil {
		return err
	}
	return nil
}

//getPath returns full url for a request
// if an id is provided the path will include
//otherwise not
func (c *Client) getPath(uri string, id string) (string, error) {

	var rurl string

	if len(id) == 0 {
		rurl = fmt.Sprintf("%s/%s", c.svcurl, uri)
	} else {
		rurl = fmt.Sprintf("%s/%s/%s", c.svcurl, uri, id)
	}

	var err error
	_, err = url.Parse(rurl)
	if err != nil {
		return "", fmt.Errorf("InvalidEndpointURL [%s]", rurl)
	}
	fmt.Printf("Built endpoint: %s\n", rurl)
	return rurl, nil
}

//getMessage accepts an interface to convert
//to []byte to be used as message sent to server
func getMessage(i interface{}) ([]byte, error) {

	t, err := json.Marshal(&i)
	if err != nil {
		return nil, fmt.Errorf("Unable to create json message. %s", err)
	}
	if t == nil || len(t) < 1 {
		return nil, fmt.Errorf("Unable to create json message. No message content")
	}
	return t, nil
}

//send sends an *http.Request and returns a byte array and error
//The byte array will be nil if there is a non nil error

func (c *Client) send(req *http.Request) ([]byte, error) {

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	//get the response stuff
	code, data, err := parseResponse(resp)
	log.Printf("code [%d], data [%d], Error [%v]", code, len(data), err)
	//check for errors in response
	if err != nil {
		return nil, err
	}
	//check for errors using response code and contents of resp.body
	err = handleError(code, data)
	log.Printf("Http Send: code [%d], data [%d], Error [%v]", code, len(data), err)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//ErrResponse is struct used to parse response body
//for json error message. Note: This is dependent of service
//error handling implementation. The following is based on a
//default Spring Restful service.
type ErrResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

//handleError takes in response code and response body in bytes
//and determines if there is a 'soft' error in what the service
//returned. Basically determines if code is something that can be
// linked with an error. Then checks the message response for json
//error message. Returns nil error if the codes and response message
//checks out, otherwise returns error per error formatting in the
//message itself.
func handleError(code int, data []byte) (err error) {

	if code == 200 {
		return nil
	}
	//no special handling based on body
	if len(data) < 1 {
		return nil
	}

	e := ErrResponse{}
	json.Unmarshal([]byte(data), &e)
	log.Printf("Error [%v], Message [%s] --> %d", e.Error, e.Message, len(e.Error))

	if len(e.Error) > 0 {
		return fmt.Errorf("%d - %s : %s", code, e.Error, e.Message)
	}
	return nil
}

//handleResponse pulls status code and body content from http response
//if no body returns zero byte array
func parseResponse(response *http.Response) (code int, data []byte, err error) {

	//pull out status code
	code = response.StatusCode

	//get response body if we have one and
	//defer the body closure
	if response.Body != nil {
		defer response.Body.Close()
		data, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return code, nil, fmt.Errorf("Unable to read response body")
		}
		return code, data, nil
	}
	return code, nil, nil
}
