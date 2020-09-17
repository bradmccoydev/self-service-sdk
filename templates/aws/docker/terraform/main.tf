##################################################################
# 
# Terraform script to manage creation/update of Lambda.
# 
##################################################################

###
# AWS provider
###
provider "aws" {
   region = var.service_aws_region
}


###
# Lambda source zip
###
data "archive_file" "lambda_zip" {
   type        = "zip"
   source_dir  = var.service_zip_input
   output_path = var.service_zip_output
}


###
# Lambda definition
###
resource "aws_lambda_function" "lambda_service" {
   function_name    = var.service_name
   description      = var.service_desc
   role             = var.service_role
   handler          = var.service_handler
   runtime          = var.service_runtime
   memory_size      = var.service_memory
   timeout          = var.service_timeout
   filename         = data.archive_file.lambda_zip.output_path
   source_code_hash = data.archive_file.lambda_zip.output_base64sha256
}


###
# Cloudwatch log definition
###
resource "aws_cloudwatch_log_group" "lambda_service_logs" {
   name              = "/aws/lambda/${aws_lambda_function.lambda_service.function_name}"
   retention_in_days = var.service_log_retention
}