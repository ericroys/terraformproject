package widget

import (
	"github.com/ericroys/terraform-provider-widget/widget/client"
)

/*Config is a configuration structure for terraform provider */
type Config struct {
	Username   string
	Password   string
	ServiceURL string
	MaxRetries int
}

//GetClient returns an initialized api client
func (c *Config) GetClient() (interface{}, error) {
	client, err := client.NewClient(c.ServiceURL)
	//client := Client{}.NewClient(c)
	return client, err
}
