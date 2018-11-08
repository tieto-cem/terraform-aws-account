data "aws_lambda_invocation" "invoce_workmail" {
  function_name = "${var.workmail_function_name}"
  input = <<JSON
{
  "action": "create-group",
  "region": "${var.region}",
  "organizationAlias": "${var.workmail_organization_alias}",
  "groupEmail": "${var.account_email}",
  "groupName": "${replace(var.account_name, " ", "")}",
  "userEmails": ${jsonencode(var.receivers)}
}
JSON
}

resource "aws_organizations_account" "account" {
  name  = "${var.account_name}"
  email = "${var.account_email}"
  iam_user_access_to_billing = "${var.access_to_billing}"
}

module "admin_role" {
  source = "./modules/admin_role"
  region = "${var.region}"
  source_account_id = "671597299301"
  account_id = "${aws_organizations_account.account.id}"
}

module "identity_assume_role_group" {
  source = "./modules/identity_assume_role_group"
  account_name = "${var.account_name}"
  region = "${var.region}"
  source_account_id = "${var.identity_source_account_id}"
  assume_role_account_id = "${aws_organizations_account.account.id}"
}
