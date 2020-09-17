##################################################################
# 
# Terraform variables
# 
##################################################################

#------------------------------------------------------------
# The following variables are for the lambda definition
#

variable "service_aws_region" {
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

variable "service_zip_input" {
  description = "The directory containing the service binary"
  type = string
}

variable "service_zip_output" {
  description = "The name & path of the zip file to be created"
  type = string
}


#------------------------------------------------------------
# The following variables are for the terraform backend
#

#variable "use_remote_terraform_state" {
#  description = "Whether to use remote storage for terraform state"
#  type = bool
#}

#variable "terraform_state_s3_bucket" {
#  description = "The name of the S3 bucket to use to store the terraform state file"
#  type = string
#}

#variable "terraform_state_s3_key" {
#  description = "The key to use to access the terraform state file"
#  type = string
#}

#variable "terraform_state_s3_region" {
#  description = "The region for the S3 bucket"
#  type = string
#}

#variable "terraform_state_dynamo_table" {
#  description = "The name of the Dynamo DB table to use for managing the terraform state locks"
#  type = string
#}
