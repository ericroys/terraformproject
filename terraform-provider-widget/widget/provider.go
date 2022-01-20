package widget

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "username",
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "password",
			},

			"service_url": {
				Type:        schema.TypeString,
				Optional:    false,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_BASEURL", nil),
				Description: "service_url",
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     25,
				Description: "max_retries",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			/*nothing defined for data sources*/
			//			"test":                    	test(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"widget_widget": resourceWidget(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		ServiceURL: d.Get("service_url").(string),
		MaxRetries: d.Get("max_retries").(int),
	}
	c, err := config.GetClient()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	return c, err
}
