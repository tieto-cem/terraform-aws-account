data "aws_iam_policy_document" "assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "bucket_access" {
  statement {
    actions = [
      "s3:GetBucketLocation",
      "s3:ListAllMyBuckets"
    ]
    resources = ["arn:aws:s3:::*"]
  }
  statement {
    actions = [
      "s3:ListBucket",
      "s3:GetObject"
    ]
    resources = [
      "arn:aws:s3:::${var.bucket}",
      "arn:aws:s3:::${var.bucket}/*"
    ]
  }
}

resource "aws_iam_role" "lambda_workmail" {
  name = "lambda_workmail"
  assume_role_policy = "${data.aws_iam_policy_document.assume_role.json}"
}

resource "aws_iam_instance_profile" "lambda_instance_profile" {
  name = "lambda_instance_profile"
  role = "lambda_workmail"
}

resource "aws_iam_role_policy" "lambda_role_policy" {
  name   = "lambda_role_policy"
  role   = "${aws_iam_role.lambda_workmail.id}"
  policy = "${data.aws_iam_policy_document.bucket_access.json}"
}

resource "aws_iam_role_policy_attachment" "AmazonWorkMailFullAccess" {
  role = "${aws_iam_role.lambda_workmail.id}"
  policy_arn = "arn:aws:iam::aws:policy/AmazonWorkMailFullAccess"
}

resource "aws_lambda_function" "aws_account_workmail" {
  function_name = "aws_account_workmail"
  s3_bucket = "${var.bucket}"
  s3_key    = "${var.file}"
  role      = "${aws_iam_role.lambda_workmail.arn}"
  handler   = "workmail"
  runtime   = "go1.x"
}
