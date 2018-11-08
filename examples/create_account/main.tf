provider "aws" {
  region = "eu-west-1"
}

module "aws_account" {
  source = "../../"
  region = "eu-west-1"
  account_name = "sandbox"
  account_email = "sandbox@example.com"
  workmail_function_name = "aws_account_workmail"
  workmail_organization_alias = "example"
  receivers = ["aws-accounts@example.com"]
  identity_source_account_id = "012345678901"
}
