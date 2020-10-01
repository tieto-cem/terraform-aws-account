provider "aws" {
  alias  = "source_account"
  region = var.region
  assume_role {
    role_arn = "arn:aws:iam::${var.source_account_id}:role/OrganizationAccountAccessRole"
  }
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    actions   = ["sts:AssumeRole"]
    resources = ["arn:aws:iam::${var.assume_role_account_id}:role/Admin"]
  }
}

resource "aws_iam_group" "assume_role_group" {
  provider = aws.source_account
  name     = "${replace(var.account_name, " ", "")}Admin"
}

resource "aws_iam_policy" "assume_role_policy" {
  provider    = aws.source_account
  name        = "CanAssumeAdminRoleIn${replace(var.account_name, " ", "")}"
  description = "Can assume admin role in ${var.account_name}"
  policy      = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_group_policy_attachment" "assume_admin_role" {
  provider   = aws.source_account
  group      = aws_iam_group.assume_role_group.name
  policy_arn = aws_iam_policy.assume_role_policy.arn
}
