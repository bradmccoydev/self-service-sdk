#!/bin/bash
#############################################################
#
# The following variables are used to configure the Lambda
# service in AWS. 
#
#############################################################

###
# Run Mode
###
#export RUN_MODE="delete"


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
# Service memory (MB in multiples of 64)
###
export SERVICE_MEMORY="128"


###
# Service storage
###
export SERVICE_STORAGE="s3://my_s3_bucket/some/path"


###
# Service storage
###
export SERVICE_RUNTIME="go1.x"


###
# Service Cloudwatch log retention period (days)
###
export SERVICE_LOG_RETENTION="7"


###
# Service role action
###
export SERVICE_ROLE_ACTION="use"


###
# Service role name
###
export SERVICE_ROLE_NAME="false"
