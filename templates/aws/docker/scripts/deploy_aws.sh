#!/bin/bash
##################################################################
#
# Script to manage build and publish of GoLang Self Service Lambda
#
##################################################################

###
# Directory Variables
###
export DIR_BASE="/tmp"
export DIR_BUILD="${DIR_BASE}/build"
export DIR_SOURCE="${DIR_BASE}/source"
export DIR_TERRAFORM="${DIR_BASE}/terraform"
export DIR_WORK="${DIR_BASE}/work"
export DIR_WORK_LOG="${DIR_WORK}/log"
export DIR_WORK_OUT="${DIR_WORK}/outputs"
export DIR_ZIP="${DIR_BASE}/zip"

###
# File Name Variables
###
export FILE_NAME_LOG="deploy_aws.log"
export FILE_NAME_SERVICE_BINARY="main"
export FILE_NAME_SERVICE_ZIP="main.zip"
export FILE_NAME_TF_BACKEND="backend.tf"
export FILE_NAME_TF_PLAN_OUT="tfapply.out"
export FILE_NAME_TFVARS="variables.tfvars"
export FILE_NAME_USER_INPUTS="inputs.sh"

###
# File Path Variables
###
export FILE_PATH_LOG="${DIR_WORK_LOG}/${FILE_NAME_LOG}"
export FILE_PATH_SERVICE_BINARY="${DIR_BUILD}/${FILE_NAME_SERVICE_BINARY}"
export FILE_PATH_SERVICE_ZIP="${DIR_ZIP}/${FILE_NAME_SERVICE_ZIP}"
export FILE_PATH_TF_BACKEND="${DIR_TERRAFORM}/${FILE_NAME_TF_BACKEND}"
export FILE_PATH_TF_PLAN_OUT="${DIR_TERRAFORM}/${FILE_NAME_TF_PLAN_OUT}"
export FILE_PATH_TFVARS="${DIR_TERRAFORM}/${FILE_NAME_TFVARS}"
export FILE_PATH_USER_INPUTS="${DIR_WORK}/${FILE_NAME_USER_INPUTS}"

###
# AWS Dynamo DB Variables
###
export AWS_DYNAMO_DB_SERVICE_TABLE="service"

###
# AWS Lambda Variables
###
export AWS_LAMBDA_HANDLER="tmp/build/lambda/main"



#------------------------------------------------------------
# Nothing should be modified below this point!!!
#

###
#
# Function to show usage
#
###
do_usage() {

   scriptName=`basename "$0"`
   echo ""
   echo "Usage: ${scriptName} [OPTION]..."
   echo ""
   echo "This script is used to build & deploy a GoLang Lambda function."
   echo "If the Lambda does not exist it will be created. If it does"
   echo "already exist then it will be updated."
   echo ""
   echo "Options:"
   echo "   -v      Verbose logging"
   echo "   -h      Displays this help and then exits"
   echo ""
}


###
#
# Function to log debug info to a log file
#
###
log_debug() {

   # This function needs one argument:
   #    => $1 is the log message

   TS=`date +"%Y-%m-%d %H:%M:%S"`
   echo "${TS} ${1}" >> ${FILE_PATH_LOG}
}


###
#
# Function to log info to stdout
#
###
log_info() {

   # This function needs two arguments:
   #    => $1 is the indentation level
   #    => $2 is log message

   case "${1}" in
      0)
         printf "%s\n" '' "${2}"
	      ;;
      1)
         printf "%3s%s\n" '' "${2}"
	      ;;
      2)
         printf "%6s%s\n" '' "${2}"
         ;;
      3)
         printf "%9s%s\n" '' "${2}"
         ;;
      4)
         printf "%12s%s\n" '' "${2}"
         ;;
      5)
         printf "%15s%s\n" '' "${2}"
         ;;
      *)
         printf "%s\n" "${2}"
         ;;
   esac

   # Do we need to write to the log also?
   if [[ ${SCRIPT_VERBOSE,,} == "true" ]]; then
      log_debug "${2}"
   fi
}


###
#
# Function to perform initial setup
#
###
do_setup() {

   # Setup build directory
   if [[ -d ${DIR_BUILD} ]]; then
      rm -rf ${DIR_BUILD}
   fi
   mkdir -p ${DIR_BUILD}

   # Setup log directory
   if [[ -d ${DIR_WORK_LOG} ]]; then
      rm -rf ${DIR_WORK_LOG}
   fi
   mkdir -p ${DIR_WORK_LOG}

   # Setup output directory
   if [[ -d ${DIR_WORK_OUT} ]]; then
      rm -rf ${DIR_WORK_OUT}
   fi
   mkdir -p ${DIR_WORK_OUT}
}


###
#
# Function to load user provided values
#
###
do_load_values() {

   # Log start
   log_info 1 "Importing user variables..."

   # Read the file in
   . "${FILE_PATH_USER_INPUTS}"

   # Verbose logging
   if [[ ${SCRIPT_VERBOSE,,} == "true" ]]; then
      log_debug "The following SCRIPT values were loaded:"
      env | grep ^SCRIPT_ | while read line
      do
         log_debug "${line}"
      done
      log_debug "The following SERVICE values were loaded:"
      env | grep ^SERVICE_ | while read line
      do
         log_debug "${line}"
      done
      log_debug "The following TERRAFORM values were loaded:"
      env | grep ^TERRAFORM_ | while read line
      do
         log_debug "${line}"
      done
   fi
}


###
#
# Function to check if variable is set
#
###
do_check_variable_set() {

   # This function needs one argument:
   #    => $1 is the variable name

   if [[ ${SCRIPT_VERBOSE,,} == "true" ]]; then
      log_debug "Checking for variable: ${1}"
   fi
   if [[ -z ${!1} ]]; then
      log_info 2 "*** FAILED *** ${1} is not set"
      log_info 2 ""
      exit 1;
   fi  
}


###
#
# Function to check if variable is a positive integer
#
###
do_check_variable_int() {

   # This function needs one argument:
   #    => $1 is the variable name

   if [[ ${SCRIPT_VERBOSE,,} == "true" ]]; then
      log_debug "Checking that variable: ${1} is a positive integer"
   fi

   REGEX_PATTERN="^[0-9]+$"
   if [[ ! ${!1} =~ ${REGEX_PATTERN} ]]; then
      log_info 2 "*** FAILED *** Variable ${1} value ${!1} is not a positive integer"
      log_info 2 ""
      exit 1;
   fi  
}


###
#
# Function to perform basic sanity checks
#
###
do_sanity_checks() {

   # Log start
   log_info 1 "Performing sanity checks..."

   # Check we have the user variables we need
   do_check_variable_set SERVICE_NAME
   do_check_variable_set SERVICE_DESC
   do_check_variable_set SERVICE_LOG_RETENTION
   do_check_variable_set SERVICE_MEMORY
   do_check_variable_set SERVICE_TIMEOUT
   do_check_variable_set SERVICE_STORAGE
   do_check_variable_set SERVICE_REGION
   do_check_variable_set SERVICE_RUNTIME
   do_check_variable_set SERVICE_ROLE_ACTION
   do_check_variable_set SERVICE_ROLE_NAME

   # Check that the following variables are positive integers
   do_check_variable_int SERVICE_LOG_RETENTION
   do_check_variable_int SERVICE_MEMORY
   do_check_variable_int SERVICE_TIMEOUT

   # Check terraform state variables are set if remote state requested
   if [[ ${TERRAFORM_REMOTE_STATE,,} == "true" ]]; then
      do_check_variable_set TERRAFORM_STATE_S3_BUCKET_NAME
      do_check_variable_set TERRAFORM_STATE_S3_KEY
      do_check_variable_set TERRAFORM_STATE_S3_REGION
      do_check_variable_set TERRAFORM_STATE_DYNAMODB_TABLE_NAME
   fi

   # Check source directory exists
   if [[ ${SCRIPT_VERBOSE,,} == "true" ]]; then
      log_debug "Checking that source directory: ${DIR_SOURCE} exists"
   fi
   if [[ ! -d ${DIR_SOURCE} ]]; then
      log_info 2 "*** FAILED *** The source directory ${DIR_SOURCE} does not exist"
      log_info 2 ""
      exit 1;
   fi

   # Is Terraform verbose enabled?
   if [[ ${TERRAFORM_VERBOSE,,} == "true" ]]; then
      if [[ ${SCRIPT_VERBOSE,,} == "true" ]]; then
         log_debug "Terraform verbose mode enabled. Setting TF_LOG to DEBUG"
      fi
      export TF_LOG=DEBUG
   fi
}


###
#
# Function to download any required dependencies
#
###
do_get_dependencies() {

   # Log start
   log_info 1 "Downloading dependencies..."

   # Download dependencies
   cd ${DIR_SOURCE}
   if [[ ${SCRIPT_VERBOSE,,} == "true" ]]; then
      log_debug ""
      go get ./... >> ${FILE_PATH_LOG} 2>&1
      RESULT=${?}
      log_debug ""
   else
      go get ./... > /dev/null 2>&1
      RESULT=${?}
   fi
   if [[ ${RESULT} -ne 0 ]]; then
      log_info 2 "*** FAILED *** ERROR reported by go get"
      log_info 2 ""
      exit 1;
   fi
}


###
#
# Function to perform build of the binary
#
###
do_build() {

   # Log start
   log_info 1 "Performing build of source..."

   # Verbose logging
   if [[ ${SCRIPT_VERBOSE,,} == "true" ]]; then
      log_debug "Found source files:"
      for FILE in `ls ${DIR_SOURCE}/*.go`
      do
         log_debug "${FILE}"
      done
   fi
 
   # Build the binary
   go build -o ${FILE_PATH_SERVICE_BINARY} 
   if [[ ${?} -ne 0 ]]; then
      log_info 2 "*** FAILED *** ERROR reported by go build"
      log_info 2 ""
      exit 1;
   fi

   # Dump a copy?
   if [[ ${SCRIPT_DUMP_ARTEFACTS,,} == "true" ]]; then

      # Copy the binary
      cp ${FILE_PATH_SERVICE_BINARY} ${DIR_WORK_OUT}
   fi
}


###
#
# Function to append an entry to tfvars
#
###
do_append_tfvar() {

   # This function needs two arguments:
   #    => $1 is the key
   #    => $2 is log value

   # Verbose logging
   if [[ ${SCRIPT_VERBOSE,,} == "true" ]]; then
      log_debug "Appending to TFVARS variable: ${1} with value: ${2}"
   fi

   # Ok, append it
   echo "${1} = \"${2}\"" >> ${FILE_PATH_TFVARS}
}


###
#
# Function to create a TFVARS file
#
###
do_create_tfvars() {

   # Log start
   log_info 1 "Creating tfvars file..."

   # Delete file if it already exists
   if [[ -f ${FILE_PATH_TFVARS} ]]; then
      rm -rf ${FILE_PATH_TFVARS}
   fi

   # Add the user variables
   do_append_tfvar "service_name"          "${SERVICE_NAME}"
   do_append_tfvar "service_desc"          "${SERVICE_DESC}"
   do_append_tfvar "service_memory"        "${SERVICE_MEMORY}"
   do_append_tfvar "service_timeout"       "${SERVICE_TIMEOUT}"
   do_append_tfvar "service_runtime"       "${SERVICE_RUNTIME}"
   do_append_tfvar "service_role"          "${SERVICE_ROLE_NAME}" 
   do_append_tfvar "service_aws_region"    "${SERVICE_REGION}"
   do_append_tfvar "service_log_retention" "${SERVICE_LOG_RETENTION}"

   # Add the other bits we need for the service
   do_append_tfvar "service_handler"       "${AWS_LAMBDA_HANDLER}"
   do_append_tfvar "service_zip_input"     "${DIR_BUILD}" 
   do_append_tfvar "service_zip_output"    "${FILE_PATH_SERVICE_ZIP}"
}


###
#
# Function to append an entry to backend.tf
#
###
do_append_backend() {

   # This function needs two arguments:
   #    => $1 is the value to append

   # Verbose logging
   if [[ ${SCRIPT_VERBOSE,,} == "true" ]]; then
      log_debug "Appending to backend file ${FILE_NAME_TF_BACKEND} the value: ${1}"
   fi

   # Ok, append it
   echo "${1}" >> ${FILE_PATH_TF_BACKEND}
}


###
#
# Function to create a backend.tf file
#
###
do_create_tf_backend() {

   # Log start
   log_info 1 "Creating ${FILE_NAME_TF_BACKEND} file..."

   # Delete file if it already exists
   if [[ -f ${FILE_PATH_TF_BACKEND} ]]; then
      rm -rf ${FILE_PATH_TF_BACKEND}
   fi

   # Remote or local backend?
   if [[ ${TERRAFORM_REMOTE_STATE,,} == "true" ]]; then
      if [[ ${SCRIPT_VERBOSE,,} == "true" ]]; then
         log_debug "Using remote AWS backend for terraform state..."
      fi
      do_append_backend "terraform {" 
      do_append_backend "   backend \"s3\" {"  
      do_append_backend "      bucket = \"${TERRAFORM_STATE_S3_BUCKET_NAME}\""
      do_append_backend "      key = \"${TERRAFORM_STATE_S3_KEY}\"" 
      do_append_backend "      region = \"${TERRAFORM_STATE_S3_REGION}\"" 
      do_append_backend "      dynamodb_table = \"${TERRAFORM_STATE_DYNAMODB_TABLE_NAME}\""
      do_append_backend "      encrypt = true" 
      do_append_backend "   }" 
      do_append_backend "}" 
   else
      if [[ ${SCRIPT_VERBOSE,,} == "true" ]]; then
         log_debug "Using local backend for terraform state..."
      fi
      do_append_backend "terraform {" 
      do_append_backend "   backend \"local\" {"  
      do_append_backend "   }" 
      do_append_backend "}" 
   fi

   # Dump a copy?
   if [[ ${SCRIPT_DUMP_ARTEFACTS,,} == "true" ]]; then

      # Copy the terraform files
      cp ${DIR_TERRAFORM}/* ${DIR_WORK_OUT}
   fi
}


###
#
# Function to run Terraform init
#
###
do_perform_tfinit() {

   # Log start
   log_info 1 "Performing terraform init..."

   # Run init
   cd ${DIR_TERRAFORM}
   if [[ ${SCRIPT_VERBOSE,,} == "true" ]]; then
      log_debug ""
      terraform init >> ${FILE_PATH_LOG} 2>&1
      RESULT=${?}
      log_debug ""
   else
      terraform init > /dev/null
      RESULT=${?}
   fi
   if [[ ${RESULT} -ne 0 ]]; then
      log_info 2 "*** FAILED *** ERROR reported by terraform init"
      log_info 2 ""
      exit 1;
   fi
}


###
#
# Function to run Terraform plan
#
###
do_perform_tfplan() {

   # Log start
   log_info 1 "Performing terraform plan..."

   # Run plan
   cd ${DIR_TERRAFORM}
   if [[ ${SCRIPT_VERBOSE,,} == "true" ]]; then
      log_debug ""
      terraform plan -input=false -out ${FILE_PATH_TF_PLAN_OUT} -var-file ${FILE_PATH_TFVARS} >> ${FILE_PATH_LOG} 2>&1
      RESULT=${?}
      log_debug ""
   else
      terraform plan -input=false -out ${FILE_PATH_TF_PLAN_OUT} -var-file ${FILE_PATH_TFVARS} > /dev/null
      RESULT=${?}
   fi
   if [[ ${RESULT} -ne 0 ]]; then
      log_info 2 "*** FAILED *** ERROR reported by terraform plan"
      log_info 2 ""
      exit 1;
   fi
}


###
#
# Function to run Terraform apply
#
###
do_perform_tfapply() {

   # Log start
   log_info 1 "Performing terraform apply..."

   # Run plan
   cd ${DIR_TERRAFORM}
   if [[ ${SCRIPT_VERBOSE,,} == "true" ]]; then
      log_debug ""
      terraform apply -input=false -auto-approve ${FILE_PATH_TF_PLAN_OUT} >> ${FILE_PATH_LOG} 2>&1
      RESULT=${?}
      log_debug ""
   else
      terraform apply -input=false -auto-approve ${FILE_PATH_TF_PLAN_OUT} > /dev/null
      RESULT=${?}
   fi
   if [[ ${RESULT} -ne 0 ]]; then
      log_info 2 "*** FAILED *** ERROR reported by terraform apply"
      log_info 2 ""
      exit 1;
   fi
}



#################################################
#
# Main entrypoint of script
#
#################################################
export MODE="POST"
export SCRIPT_VERBOSE="FALSE"


# Check that we received valid options
options=':h'
while getopts $options option
do
   case ${option} in
      h  ) do_usage; exit 0 ;;
      \? ) echo ""; echo "ERROR: Unknown option: $OPTARG" 1>&2; exit 1;;
   esac
done


# Setup
do_setup

# Start processing
log_info 0 ""
log_info 0 "###################################################"
log_info 0 "Start of processing..."

# Load the values provided by the user
do_load_values

# Perform sanity checks
do_sanity_checks

# If we are not in delete mode
if [[ ${MODE} != "DELETE" ]]; then

   # Get any build dependencies
   do_get_dependencies

   # Build the binary
   do_build

   # Create TFVARS file
   do_create_tfvars

   # Create backend state file
   do_create_tf_backend

   # Perform terraform init
   do_perform_tfinit

   # Perform terraform plan
   do_perform_tfplan

   # Perform terraform apply
   do_perform_tfapply

fi

# Log finish
log_info 0 "End of processing..."
log_info 0 "###################################################"
