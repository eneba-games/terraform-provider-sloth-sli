package sli

import (
	"context"

	"github.com/HelisLT/terraform-provider-sloth-sli/generator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	paramSlothPath = "sloth_path"
)

func Provider(configureContextFunc schema.ConfigureContextFunc) *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			paramSlothPath: {
				Type:     schema.TypeString,
				Required: true,
				Description: `Path to "sloth" executable.
How to install it can be found here: https://sloth.dev/introduction/install/`,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sli": dataSourceSLI(),
		},
		ConfigureContextFunc: configureContextFunc,
	}
}

// GeneratorConfigureFunc is a default provider's ConfigureContextFunc which configures SLI generator
// with given sloth executable
func GeneratorConfigureFunc(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	slothPath := d.Get(paramSlothPath).(string)

	return generator.NewSLI(slothPath), diags
}
