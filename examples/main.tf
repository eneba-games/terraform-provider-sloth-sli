terraform {
  required_providers {
    sli = {
      source  = "eneba-games/sloth-sli"
      version = "~> 0.0.1"
    }
  }
}

provider "sli" {
  sloth_path = "/usr/local/bin/sloth"
}

data "sli" "http_latency_config" {
  sli_config = <<EOT
version: "prometheus/v1"
service: "myservice"
labels:
  owner: "myteam"
  repo: "myorg/myservice"
  tier: "2"
slos:
  # We allow failing (5xx and 429) 1 request every 1000 requests (99.9%).
  - name: "requests-availability"
    objective: 99.9
    description: "Common SLO based on availability for HTTP request responses."
    sli:
      events:
        error_query: sum(rate(http_request_duration_seconds_count{job="myservice",code=~"(5..|429)"}[{{.window}}]))
        total_query: sum(rate(http_request_duration_seconds_count{job="myservice"}[{{.window}}]))
    alerting:
      name: MyServiceHighErrorRate
      labels:
        category: "availability"
      annotations:
        # Overwrite default Sloth SLO alert summmary on ticket and page alerts.
        summary: "High error rate on 'myservice' requests responses"
      page_alert:
        labels:
          severity: pageteam
          routing_key: myteam
      ticket_alert:
        labels:
          severity: "slack"
          slack_channel: "#alerts-myteam"
EOT
}

output "http_latency_rules" {
  value = data.sli.http_latency_config.rendered
}
