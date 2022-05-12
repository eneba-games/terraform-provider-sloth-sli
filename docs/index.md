---
page_title: "Provider: Sloth SLI"
subcategory: ""
description: |-
  Terraform provider for generating prometheus recording and alerting rules by [sloth](https://sloth.dev/) SLI config.
---

# Sloth SLI Provider

-> Visit the [sloth](https://sloth.dev/) to learn more about the Sloth tool.

The Sloth SLI provider is used to generate prometheus recording and alerting rules by given sloth SLI config.
Generated rules can be later uploaded to desired location, e.g. prometheus server, S3 bucket, etc.
It depends on how prometheus is being provisioned.

## Example Usage

Provider requires path to `sloth` executable

```terraform
provider "sli" {
  sloth_path = "/usr/local/bin/sloth"
}
```

## Schema

### Required

- **sloth_path** (String, Required) Path to `sloth` executable
