variable "region" {
  description = "AWS region to use for the lambda function"
  default     = "eu-west-1"
}

variable "account_name" {
  description = "AWS account name"
}

variable "assume_role_account_id" {
  description = "Account to assume role"
}

variable "source_account_id" {
  description = "Account id for the identity account"
}

