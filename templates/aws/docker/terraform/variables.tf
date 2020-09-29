##################################################################
# 
# Terraform variables
# 
##################################################################

###
# Common variables
###
variable "common_aws_region" {
  description = "AWS region"
  type = string
}


###
# Lambda function variables
###

variable "lambda_function_desc" {
  description = "The description for the Lambda function"
  type = string
}

variable "lambda_function_handler" {
  description = "The service handler path for the lambda function"
  type = string
}

variable "lambda_function_memory" {
  description = "The memory limit for the lambda function"
  type = number
}

variable "lambda_function_name" {
  description = "The name for the lambda function"
  type = string
}

variable "lambda_function_role" {
  description = "The IAM role for the lambda function"
  type = string
}

variable "lambda_function_runtime" {
  description = "The runtime type for the lambda function"
  type = string
}

variable "lambda_function_tags" {
  description = "Tags to apply to the lambda function"
  type = map(string)
  default = {}
}

variable "lambda_function_timeout" {
  description = "The timeout for the lambda function"
  type = number
}

variable "lambda_function_vars" {
  description = "Environment variables to apply to the lambda function"
  type = map(string)
  default = {}
}


###
# Lambda source variables
###

variable "lambda_source_zip_input" {
  description = "The directory containing the service binary"
  type = string
}

variable "lambda_source_zip_output" {
  description = "The name & path of the zip file to be created"
  type = string
}


###
# Cloudwatch log variables
###

variable "cloudwatch_log_retention" {
  description = "The retention period for the lambda function logs"
  type = number
}
