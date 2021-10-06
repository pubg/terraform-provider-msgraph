---
layout: "msgraph"
subcategory: "ADGroup"
page_title: "MsGraph: msgraph_groups"
description: |-
  Get AWS CloudTrail Service Account ID for storing trail data in S3.
---

# Data Source: msgraph_groups

Use this data source to get the ARN of a certificate in AWS Certificate
Manager (ACM), you can reference
it by domain without having to hard code the ARNs as input.

## Example Usage

```terraform
# Find a certificate that is issued
data "aws_acm_certificate" "issued" {
  domain   = "tf.example.com"
  statuses = ["ISSUED"]
}

# Find a certificate issued by (not imported into) ACM
data "aws_acm_certificate" "amazon_issued" {
  domain      = "tf.example.com"
  types       = ["AMAZON_ISSUED"]
  most_recent = true
}

# Find a RSA 4096 bit certificate
data "aws_acm_certificate" "rsa_4096" {
  domain    = "tf.example.com"
  key_types = ["RSA_4096"]
}
```
