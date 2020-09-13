##################################################################
# 
# Terraform variables
# 
##################################################################

variable "aws_region" {
  description = "AWS region"
}

variable "service_name" {
  description = "Lambda service name"
}

variable "service_desc" {
  description = "Lambda service description"
}

variable "service_handler" {
  description = "Lambda service handler path"
}

variable "service_runtime" {
  description = "Lambda service runtime type"
}

variable "service_role" {
  description = "Lambda service IAM role"
}

variable "service_memory" {
  description = "Lambda service memory limit"
}

variable "service_timeout" {
  description = "Lambda service timeout"
}

variable "service_log_retention" {
  description = "Lambda service Cloudwatch log retention period"
}
