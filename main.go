package main

import (
	"github.com/HelisLT/terraform-provider-sloth-sli/sli"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return sli.Provider(sli.GeneratorConfigureFunc)
		},
	})
}
