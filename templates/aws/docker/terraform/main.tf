###
# AWS provider
###
provider "aws" {
   region = var.common_aws_region
}


###
# Lambda source zip
###
data "archive_file" "lambda_zip" {
   type        = "zip"
   source_dir  = var.lambda_source_zip_input
   output_path = var.lambda_source_zip_output
}


###
# Lambda definition
###
resource "aws_lambda_function" "lambda_service" {
   function_name    = var.lambda_function_name
   description      = var.lambda_function_desc
   role             = var.lambda_function_role
   handler          = var.lambda_function_handler
   runtime          = var.lambda_function_runtime
   memory_size      = var.lambda_function_memory
   timeout          = var.lambda_function_timeout
   tags             = var.lambda_function_tags
   filename         = data.archive_file.lambda_zip.output_path
   source_code_hash = data.archive_file.lambda_zip.output_base64sha256

   dynamic "environment" {
      for_each = length(keys(var.lambda_function_vars)) == 0 ? [] : [true]
      content {
         variables = var.lambda_function_vars
      }
   }
}


###
# Cloudwatch log definition
###
resource "aws_cloudwatch_log_group" "lambda_service_logs" {
   name              = "/aws/lambda/${aws_lambda_function.lambda_service.function_name}"
   retention_in_days = var.cloudwatch_log_retention
}
