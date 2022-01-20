package widget

import (
	"log"
	"strings"

	"github.com/ericroys/terraform-provider-widget/widget/client"

	"github.com/hashicorp/terraform/helper/schema"
)

const wName string = "name"
const wUID string = "uid"
const wID string = "id"
const wIDD string = "wid"
const wSize string = "size"
const widget string = "/widget/"

func resourceWidget() *schema.Resource {
	return &schema.Resource{
		Create: resourceWidgetCreate,
		Read:   resourceWidgetRead,
		Update: resourceWidgetUpdate,
		Delete: resourceWidgetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			wName: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			wUID: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			wIDD: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			wSize: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceWidgetCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	//map resource vals into struct
	nw := client.WidgetNew{
		Name: d.Get(wName).(string),
		Size: d.Get(wSize).(string),
	}
	//create widget
	w, err := c.CreateWidget(nw)
	if err != nil {
		return err
	}
	//set the new widget id to the resource id
	d.SetId(w.ID)

	return resourceWidgetRead(d, m)
}

func resourceWidgetRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	//get widget using resource id
	w, err := c.GetWidget(d.Id())
	if err != nil {
		return err
	}
	//set resource fields from returned widget
	d.SetId(w.ID)
	d.Set(wUID, w.UID)
	d.Set(wName, w.Name)
	d.Set(wSize, w.Size)

	return nil
}

func resourceWidgetUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	//check name and size for changes
	if d.HasChange(wName) || d.HasChange(wSize) {
		id := d.Id()
		//map resource vals into struct
		nw := client.WidgetNew{
			Name: d.Get(wName).(string),
			Size: d.Get(wSize).(string),
		}
		//call update with new vals
		_, err := c.UpdateWidget(id, nw)
		if err != nil {
			return err
		}
	}
	return resourceWidgetRead(d, m)
}

func resourceWidgetDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	log.Printf("[INFO] Deleting widget [%s]", d.Id())

	//delete the widget by id
	err := c.DeleteWidget(d.Id())
	//return error unless widget isn't found
	if err != nil && !strings.Contains(err.Error(), "Widget Not Found") {
		return err
	}
	//set resource id to empty
	d.SetId("")

	return nil
}
