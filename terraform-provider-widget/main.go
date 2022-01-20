package main

import (
	"github.com/ericroys/terraform-provider-widget/widget"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: widget.Provider})

}
