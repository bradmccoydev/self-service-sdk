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
# Service memory (MB in multiples of 64)
###
export SERVICE_MEMORY="128"


###
# Service storage
###
export SERVICE_STORAGE="s3://my_s3_bucket/some/path"


###
# Service runtime
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


###
# Service AWS Region
###
export SERVICE_REGION="us-west-2"


#----------------------------------------------------------#
# The follow variables control the behaviour of the 
# deployment script.
#

###
# Run Mode
###
#export SERVICE_DEPLOY_MODE="delete"


###
# Dump artefacts
###
export SERVICE_ARTEFACTS_DUMP="true"
