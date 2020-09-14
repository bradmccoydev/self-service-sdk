##################################################################
# 
# Terraform variables
# 
##################################################################

variable "aws_region" {
  description = "AWS region"
  type = string
}

variable "service_name" {
  description = "Lambda service name"
  type = string
}

variable "service_desc" {
  description = "Lambda service description"
  type = string
}

variable "service_handler" {
  description = "Lambda service handler path"
  type = string
}

variable "service_runtime" {
  description = "Lambda service runtime type"
  type = string
}

variable "service_role" {
  description = "Lambda service IAM role"
  type = string
}

variable "service_memory" {
  description = "Lambda service memory limit"
  type = number
}

variable "service_timeout" {
  description = "Lambda service timeout"
  type = number
}

variable "service_log_retention" {
  description = "Lambda service Cloudwatch log retention period"
  type = number
}

variable "zip_input" {
  description = "The directory containing the service binary"
  type = string
}

variable "zip_output" {
  description = "The name & path of the zip file to be created"
  type = string
}
