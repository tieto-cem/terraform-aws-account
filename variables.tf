variable "region" {
  description = "AWS region to use for the lambda function"
  default = "eu-west-1"
}

variable "account_name" {
  description = "AWS account name"
}

variable "account_email" {
  description = "Account email address"
}

variable "workmail_organization_alias" {
  description = "Organization alias for workmail"
}

variable "receivers" {
  description = "Workmail users part of the newly created group"
  type = "list"
}

variable "workmail_function_name" {
  description = "Workmail function name for handling creation of mail addresses"
}

variable "access_to_billing" {
  description = "Should IAM users have access to account billing information (ALLOW or DENY)"
  default = "ALLOW"
}

variable "identity_source_account_id" {
  description = "AWS identity account where users login"
}
