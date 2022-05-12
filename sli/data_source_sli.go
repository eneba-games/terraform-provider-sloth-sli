package sli

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	sliGenerator interface {
		Generate(inputFile string) (string, error)
	}
)

func dataSourceSLI() *schema.Resource {
	return &schema.Resource{
		ReadContext: sliRead,

		Schema: map[string]*schema.Schema{
			"sli_config": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rendered": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "rendered recording and alerting rules",
			},
		},
	}
}

func sliRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	g, ok := meta.(sliGenerator)
	if !ok {
		return diag.FromErr(errors.New("resource meta doesn't implement SLI sliGenerator interface"))
	}

	var diags diag.Diagnostics

	sliConfig := d.Get("sli_config").(string)
	fileName := fmt.Sprintf("%s/%s.yml", os.TempDir(), hash(sliConfig))

	defer func() {
		os.Remove(fileName)
	}()

	if err := ioutil.WriteFile(fileName, []byte(sliConfig), 0777); err != nil {
		return diag.FromErr(fmt.Errorf("failed to create SLI config file %s: %w", fileName, err))
	}

	out, err := g.Generate(fileName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("rendered", out)
	d.SetId(hash(out))

	return diags
}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))

	return hex.EncodeToString(sha[:])
}
