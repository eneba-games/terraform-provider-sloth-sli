---
page_title: "SLI Data Source - terraform-provider-sloth-sli"
subcategory: ""
description: |-
  The SLI data source allows you to retrieve generated prometheus recording and alerting rules by given sloth SLI config.
---

# Data Source `sli`

The SLI data source allows you to retrieve generated prometheus recording and alerting rules by given sloth SLI config.

## Example Usage

```terraform
data "sli" "http_latency_config" {
  sli_config = file("/path-to/http-latency-sli-config.yml")
}
```

## Argument Reference

- `sli_config` - The SLI config, [example](https://sloth.dev/examples/default/getting-started/).

## Attributes Reference

The following attributes are exported.

- `rendered` - A list of prometheus recording and alerting rules for given SLI config.
