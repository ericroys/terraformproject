package checkptclient

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ericroys/checkptclient/rest"
)

const (
	endpointLogin   = `login`
	endpointAddHost = `add-host`
	endpointPublish = `publish`
)

//APIConfig provides the construct for configuring the
//CheckPoint APIClient
type APIConfig struct {
	//User, Pwd string
	Baseurl  string
	CertPath string
	session  Session
}

//SetSessionLastPublish Allows login to an existing last
//published session for a user
func (ac *APIConfig) SetSessionLastPublish(last bool) {
	ac.session.SessContPub = last
}

//SetSessionContLast Allows login to the last session for a user.
func (ac *APIConfig) SetSessionContLast(last bool) {
	ac.session.SessCont = last
}

//NewAPIConfig creates and initializes an APIConfig object.
//Defaults to use the last session for the user
func NewAPIConfig(baseurl, user, pass, certpath string) *APIConfig {
	s := Session{
		User:        user,
		Password:    pass,
		SessTimeout: 600,
		//default to continue last session
		SessCont:    true,
		SessContPub: false,
	}
	ac := APIConfig{
		Baseurl:  baseurl,
		CertPath: certpath,
		session:  s,
	}
	return &ac
}

//APIClient is the CheckPoint API Client. All interaction with
//a Check Point service is done using methods provided by this
//client.
type APIClient struct {
	conf        *APIConfig
	httpClient  *http.Client
	sid         string
	nextRefresh time.Time
}

func (a *APIClient) getSID() error {
	n := time.Now()
	if a.sid == "" || n.After(a.nextRefresh) {
		err := a.Login()
		if err != nil {
			return err
		}
	}
	if a.sid == "" {
		return fmt.Errorf("unable to obtain the session identifier")
	}
	return nil
}

//Login logs into the Check Point service and
//returns a session identifier and session timeout
//to the client
func (a *APIClient) Login() error {
	var resp LoginResponse
	//l := a.conf.session
	uri, err := a.getPath(endpointLogin, "")
	if err != nil {
		return err
	}
	err = a.send(uri, &a.conf.session, &resp, false)
	if err != nil {
		return err
	}

	log.Printf("Sid: %s\nTimeout: %d\n", resp.Sid, resp.SessTimeout)
	a.sid = resp.Sid

	//pad in a 5 second buffer for the timeout
	a.nextRefresh = time.Now().Add(
		time.Duration(resp.SessTimeout-5) * time.Second)
	return nil
}

//CreateHost creates a Host on the CheckPoint service
func (a *APIClient) CreateHost(host Host) (Host, error) {
	var h Host
	uri, err := a.getPath(endpointAddHost, "")
	if err != nil {
		return h, err
	}

	err = a.send(uri, &host, &h, true)
	if err != nil {
		return h, err
	}
	return h, nil
}

func (a *APIClient) Publish() error {

	var msg NoMessage
	uri, err := a.getPath(endpointPublish, "")
	if err != nil {
		return err
	}

	err = a.send(uri, &msg, &msg, true)
	if err != nil {
		return err
	}
	return nil
}

func (a *APIClient) getSender(uri string, msg []byte, auth bool) (*rest.Request, error) {

	builder := rest.NewRequestBuilder(uri, a.httpClient).
		Auth(rest.AuthNoAuth{}).
		Header("Accept", "application/json").
		Message(msg).
		Method(rest.POST).
		ErrorHandler(ErrHandler{})

	//if auth flag then need to get the session id to pass in the
	//header.
	if auth {
		//make sure we have a current sid
		if err := a.getSID(); err != nil {
			return nil, err
		}
		//set the header with current sid
		builder.Header("X-chkp-sid", a.sid)
	}

	v, err := builder.Build()
	return v, err
}

//NewClient initializes and validates a new APIClient provided
//the APIConfig
func NewClient(conf *APIConfig) (*APIClient, error) {
	//validate the url is ok
	if err := validURL(conf.Baseurl); err != nil {
		return nil, err
	}

	//customize client so set timeout and connection pool stuff
	trans := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}

	//get certificate stuff
	if len(conf.CertPath) > 0 {
		p, err := buildCertPool(conf.CertPath)
		if err != nil {
			return nil, err
		}
		trans.TLSClientConfig = p
	}
	return &APIClient{
		conf: conf,
		httpClient: &http.Client{
			Transport: trans,
			Timeout:   20 * time.Second,
		},
	}, nil
}

//buildCertPool
func buildCertPool(crtpath string) (*tls.Config, error) {
	cert, err := ioutil.ReadFile(crtpath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read certificate file %s", crtpath)
	}

	//for self signed certs
	insecure := flag.Bool("insecure-ssl", false, "Accept/Ignore all server SSL certificates")
	flag.Parse()

	certPool, _ := x509.SystemCertPool()
	if certPool == nil {
		certPool = x509.NewCertPool()
	}
	certPool.AppendCertsFromPEM(cert)

	return &tls.Config{
		InsecureSkipVerify: *insecure,
		RootCAs:            certPool,
	}, nil
}

//checks if url is valid, errors of not
func validURL(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("api client requires a base url")
	}
	_, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("InvalidURL [%s]", s)
	}
	return nil
}

//getPath returns full url for a request
// if an id is provided the path will include
//otherwise not
func (a *APIClient) getPath(uri string, id string) (string, error) {

	var rurl string

	if len(id) == 0 {
		rurl = fmt.Sprintf("%s/%s", a.conf.Baseurl, uri)
	} else {
		rurl = fmt.Sprintf("%s/%s/%s", a.conf.Baseurl, uri, id)
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

func getResponse(data []byte, i interface{}) error {
	err := json.Unmarshal(data, &i)
	if err != nil {
		return fmt.Errorf("failed to transform response message. %v", err)
	}
	return nil
}

func (a *APIClient) send(url string, msg interface{}, resp interface{}, auth bool) error {

	//build message
	toMsg, err := getMessage(&msg)
	fmt.Println(string(toMsg))
	if err != nil {
		return err
	}
	s, err := a.getSender(url, toMsg, auth)
	if err != nil {
		return err
	}

	data, err := s.Send()
	if err != nil {
		return err
	}
	//log.Printf("data before response trans: %s", string(data))
	err = getResponse(data, &resp)
	if err != nil {
		return err
	}
	return nil
}
