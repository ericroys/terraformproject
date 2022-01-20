package widget

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ericroys/terraform-provider-widget/widget/client"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccWidget_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWidgetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWidget(t),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWidgetExists("widget_widget.test", t),
					resource.TestCheckResourceAttr(
						"widget_widget.test", "name", "terraformTest"),
					resource.TestCheckResourceAttr(
						"widget_widget.test", "size", "momentous occasion"),
				),
			},
		},
	})
}

func testAccCheckWidgetExists(resource string, t *testing.T) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			fmt.Printf("ERR: Not found: %s\n", resource)
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		name := rs.Primary.ID
		apiClient := testAccProvider.Meta().(*client.Client)

		_, err := apiClient.GetWidget(name) //.Get(uri+name, i)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckWidgetDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "widget_widget" {
			continue
		}

		_, err := apiClient.GetWidget(rs.Primary.ID) //(uri+rs.Primary.ID, i)
		if err == nil {
			return fmt.Errorf("Widget still exists")
		}
		if !strings.Contains(err.Error(), "Widget Not Found") {
			return err
		}
	}

	return nil
}

func testAccCheckWidget(t *testing.T) string {
	t.Log("Get Widget Resource - OK")
	return fmt.Sprintf(
		`resource "widget_widget" "test"{
			name = "terraformTest"
			size = "momentous occasion"
		  }`)
}
