package widget

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"widget": testAccProvider,
	}
}

func TestAccPreCheck(t *testing.T) {

	if _, v := os.LookupEnv("SERVICE_BASEURL"); !v {
		t.Fatal("SERVICE_BASEURL must be set for acceptance tests")
	}
	v, _ := os.LookupEnv("SERVICE_BASEURL")
	t.Logf("Pre-Check - OK  -> %s", v)
}
