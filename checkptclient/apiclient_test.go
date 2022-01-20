package checkptclient

import "testing"

func client() (*APIClient, error) {

	conf := NewAPIConfig(
		"https://r8020sms.seamlessti.net/web_api/v1.3",
		"admin",
		"vpn12345",
		`c:/test/certs/chkpnt.pem`,
	)
	return NewClient(conf)
}

func TestConfig(t *testing.T) {
	if _, err := client(); err != nil {
		t.Fatal(err)
	}
}
func TestLogin(t *testing.T) {
	c, err := client()
	if err != nil {
		t.Fatal(err)
	}

	//no extra params
	err = c.Login()
	if err != nil {
		t.Fatal(err)
	}
	//t.Logf("response: %+v", r)
}

func TestCreateHost(t *testing.T) {
	c, err := client()
	if err != nil {
		t.Fatal(err)
	}

	//setup the host
	h := Host{
		Name:        "bob$suncle",
		Ipv4address: "192.168.2.145",
	}
	//no extra params
	r, err := c.CreateHost(h)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("response: %+v", r)
}

func TestPublish(t *testing.T) {
	c, err := client()
	if err != nil {
		t.Fatal(err)
	}

	//no extra params
	err = c.Publish()
	if err != nil {
		t.Fatal(err)
	}
}
