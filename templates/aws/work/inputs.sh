#!/bin/bash
#############################################################
#
# The following variables are used to configure the Lambda
# service in AWS. 
#
#############################################################

#----------------------------------------------------------#
# The follow variables define the Lambda service name
# behaviour etc.
#

###
# Service name
###
export SERVICE_NAME="myservice"


###
# Service description
###
export SERVICE_DESC="My very special microservice"


###
# Service timeout (seconds)
###
export SERVICE_TIMEOUT="30"


###
# Service Cloudwatch log retention period (days)
###
export SERVICE_LOG_RETENTION="7"


###
# Service memory (MB in multiples of 64)
###
export SERVICE_MEMORY="128"


###
# Service AWS Region
###
export SERVICE_REGION="us-west-2"


###
# Service role action
###
export SERVICE_ROLE_ACTION="use"


###
# Service role name
###
export SERVICE_ROLE_NAME="arn:aws:iam::12345:role/MyRole"


###
# Service runtime
###
export SERVICE_RUNTIME="go1.x"


###
# Service tags
###
#export SERVICE_TAGS="tagname1=123,tagname2=456"


###
# Service variables
###
#export SERVICE_VARS="var1=123,var2=456"



#----------------------------------------------------------#
# The follow variables control the behaviour of the 
# deployment script.
#

###
# Run Mode
###
#export SCRIPT_RUN_MODE="delete"


###
# Enable verbose logging
###
export SCRIPT_VERBOSE="false"


###
# Dump artefacts
###
export SCRIPT_DUMP_ARTEFACTS="true"



#----------------------------------------------------------#
# The follow variables define whether to use a remote 
# backend (S3 & DynamoDB) for the terraform state.
#

###
# Use remote state
###
export TERRAFORM_REMOTE_STATE="false"

###
# Name of the S3 bucket to use 
###
export TERRAFORM_STATE_S3_BUCKET_NAME="some_bucket"

###
# Path to state file in the S3 bucket
###
export TERRAFORM_STATE_S3_KEY="global/s3/terraform.tfstate"

###
# Region for S3 bucket
###
export TERRAFORM_STATE_S3_REGION="us-west-2"

###
# Name of the Dynamo DB table to use
###
export TERRAFORM_STATE_DYNAMODB_TABLE_NAME="tflocks"

###
# Enable Terraform verbose logging
###
export TERRAFORM_VERBOSE="false"
