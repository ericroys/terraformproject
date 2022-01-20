package checkptclient

/* All the structures used in marshal/unmarshal of json
   to and from the Check Point service
*/

//NoMessage used of sending/receiving an empty json message
type NoMessage struct{}

//LoginResponse struct for login response message
type LoginResponse struct {
	Sid         string `json:"sid"`
	Standby     bool   `json:"standby"`
	ReadOnly    bool   `json:"read-only"`
	UID         string `json:"uid"`
	SessTimeout int    `json:"session-timeout"`
}

//Session struct for defining session parameters
//used when establishing communication with the service
type Session struct {
	User         string `json:"user"`
	Password     string `json:"password"`
	Domain       string `json:"domain,omitempty"`
	SessCont     bool   `json:"continue-last-session,omitempty"`
	SessContPub  bool   `json:"enter-last-published-session,omitempty"`
	SessComments string `json:"session-comments,omitempty"`
	SessDesc     string `json:"session-descriptions,omitempty"`
	SessName     string `json:"session-name,omitempty"`
	SessTimeout  int    `json:"session-timeout,omitempty"`
}

//Host struct for definining and marshal/unmarshal of Host object
type Host struct {
	UID         string `json:"uid,omitempty"`
	Name        string `json:"name,omitempty"`
	Ipv4address string `json:"ipv4-address,omitempty"`
	Color       string `json:"color,omitempty"`
	Newname     string `json:"new-name,omitempty"`
	NatSettings `json:"nat-settings,omitempty"`
}

//NatSettings struct for defining and marshal/unmarshal of NatSettings object
type NatSettings struct {
	Hidebehind string `json:"hide-behind,omitempty"`
	Ipaddress  string `json:"ip-address,omitempty"`
	Autorule   bool   `json:"auto-rule"`
	Installon  string `json:"install-on,omitempty"`
	Method     string `json:"method,omitempty"`
}

//ErrMsgObj is an embedded message object for errors, warnings
//and blocking errors
type ErrMsgObj struct {
	Message string `json:"message"`
	Session bool   `json:"current_session,omitempty"`
}

//ErrResponse is struct defining the error response from
//the Check Point service.
type ErrResponse struct {
	Message  string      `json:"message"`
	Warnings []ErrMsgObj `json:"warnings,omitempty"`
	Errors   []ErrMsgObj `json:"errors,omitempty"`
	Blocking []ErrMsgObj `json:"blocking-errors,omitempty"`
	Code     string      `json:"code"`
}
