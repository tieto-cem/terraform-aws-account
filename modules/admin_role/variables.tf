variable "region" {
  description = "AWS region to use for the lambda function"
  default = "eu-west-1"
}

variable "account_id" {
  description = "Account number for the newly created account"
}

variable "source_account_id" {
  description = "Account number for the identity account"
}

