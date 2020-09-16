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
export DIR_WORK_OUT="${DIR_WORK}/outputs"
export DIR_ZIP="${DIR_BASE}/zip"

###
# File Name Variables
###
export FILE_NAME_SERVICE_BINARY="main"
export FILE_NAME_SERVICE_ZIP="main.zip"
export FILE_NAME_TF_PLAN_OUT="tfapply.out"
export FILE_NAME_TFVARS="variables.tfvars"
export FILE_NAME_USER_INPUTS="inputs.sh"

###
# File Path Variables
###
export FILE_PATH_SERVICE_BINARY="${DIR_BUILD}/${FILE_NAME_SERVICE_BINARY}"
export FILE_PATH_SERVICE_ZIP="${DIR_ZIP}/${FILE_NAME_SERVICE_ZIP}"
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
# Function to log indented output
#
###
log_it() {

   # This function needs two arguments:
   #    => $1 is the indentation level
   #    => $2 is log message

   # We do need to write it to stdout
   case "$1" in
   
      0)
         printf "%s\n" '' "$2"
	      ;;

      1)
         printf "%3s%s\n" '' "$2"
	      ;;

      2)
         printf "%6s%s\n" '' "$2"
         ;;

      3)
         printf "%9s%s\n" '' "$2"
         ;;

      4)
         printf "%12s%s\n" '' "$2"
         ;;

      5)
         printf "%15s%s\n" '' "$2"
         ;;

      *)
         printf "%s\n" "$2"
         ;;

   esac
}


###
#
# Function to load user provided values
#
###
do_load_values() {

   # Log start
   log_it 1 "Importing user variables..."

   # Read the file in
   . "${FILE_PATH_USER_INPUTS}"

   # Verbose logging
   if [[ ${VERBOSE} == "TRUE" ]]; then
      log_it 2 ""
      log_it 2 "The following values were loaded:"
      for value in `env | grep ^SERVICE_`
      do
         log_it 3 ${value}
      done
      log_it 2 ""
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

   if [[ ${VERBOSE} == "TRUE" ]]; then
      log_it 2 "Checking for variable: ${1}"
   fi
   if [[ -z ${!1} ]]; then
      log_it 2 "*** FAILED *** ${1} is not set"
      log_it 2 ""
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

   if [[ ${VERBOSE} == "TRUE" ]]; then
      log_it 2 "Checking that variable: ${1} is a positive integer"
   fi

   REGEX_PATTERN="^[0-9]+$"
   if [[ ! ${!1} =~ ${REGEX_PATTERN} ]]; then
      log_it 2 "*** FAILED *** Variable ${1} value ${!1} is not a positive integer"
      log_it 2 ""
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
   log_it 1 "Performing sanity checks..."

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

   # Check source directory exists
   if [[ ${VERBOSE} == "TRUE" ]]; then
      log_it 2 "Checking that source directory: ${DIR_SOURCE} exists"
   fi
   if [[ ! -d ${DIR_SOURCE} ]]; then
      log_it 2 "*** FAILED *** The source directory ${DIR_SOURCE} does not exist"
      log_it 2 ""
      exit 1;
   fi
}


###
#
# Function to download any required dependencies
#
###
do_get_dependencies() {

   # Log start
   log_it 1 "Downloading dependencies..."

   # Download dependencies
   cd ${DIR_SOURCE}
   if [[ ${VERBOSE} == "TRUE" ]]; then
      log_it 2 ""
      go get ./... 
      RESULT=${?}
      log_it 2 ""
   else
      go get ./... > /dev/null 2>&1
      RESULT=${?}
   fi
   if [[ ${RESULT} -ne 0 ]]; then
      log_it 2 "*** FAILED *** ERROR reported by go get"
      log_it 2 ""
      log_it 2 
      log_it 2 ""
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
   log_it 1 "Performing build of source..."

   # Setup
   if [[ -d ${DIR_BUILD} ]]; then
      rm -rf ${DIR_BUILD}
   fi
   mkdir -p ${DIR_BUILD}
   

   # Verbose logging
   if [[ ${VERBOSE} == "TRUE" ]]; then
      log_it 2 "Found source files:"
      for FILE in `ls ${DIR_SOURCE}/*.go`
      do
         log_it 3 "${FILE}"
      done
      log_it 2 ""
   fi
 
   # Build the binary
   go build -o ${FILE_PATH_SERVICE_BINARY} 
   if [[ ${?} -ne 0 ]]; then
      log_it 2 "*** FAILED *** ERROR reported by go build"
      log_it 2 ""
      exit 1;
   fi

   # Dump a copy?
   if [[ ${SERVICE_ARTEFACTS_DUMP,,} == "true" ]]; then

      # Check output directory exists
      if [[ ! -d ${DIR_WORK_OUT} ]]; then
         mkdir -p ${DIR_WORK_OUT}
      fi

      # Remove file if already there
      if [[ -f ${DIR_WORK_OUT}/${FILE_NAME_SERVICE_BINARY} ]]; then
         rm -rf ${DIR_WORK_OUT}/${FILE_NAME_SERVICE_BINARY}
      fi

      # Copy the binary
      cp ${FILE_PATH_SERVICE_BINARY} ${DIR_WORK_OUT}
   fi
}


# ###
# #
# # Function to perform zip of binary
# #
# ###
# do_zip() {

#    # Log start
#    log_it 1 "Creating zip"
#    if [[ -d ${DIR_ZIP} ]]; then
#       rm -rf ${DIR_ZIP}
#    fi
#    mkdir -p ${DIR_ZIP}

#    # Zip the binary
#    if [[ ${VERBOSE} == "TRUE" ]]; then
#       zip -r9 ${FILE_SERVICE_ZIP} ${FILE_SERVICE_BINARY} 
#       log_it 2 ""
#    else
#       zip -qr9 ${FILE_SERVICE_ZIP} ${FILE_SERVICE_BINARY} 
#    fi
# }


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
   if [[ ${VERBOSE} == "TRUE" ]]; then
      log_it 2 "Appending to TFVARS variable: ${1} with value: ${2}"
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
   log_it 1 "Creating terraform tfvars file..."

   # Delete file if it already exists
   if [[ -f ${FILE_PATH_TFVARS} ]]; then
      rm -rf ${FILE_PATH_TFVARS}
   fi

   # Add the user variables
   do_append_tfvar "service_name" "${SERVICE_NAME}"
   do_append_tfvar "service_desc" "${SERVICE_DESC}"
   do_append_tfvar "service_memory" "${SERVICE_MEMORY}"
   do_append_tfvar "service_timeout" "${SERVICE_TIMEOUT}"
   do_append_tfvar "service_runtime" "${SERVICE_RUNTIME}"
   do_append_tfvar "service_role" "${SERVICE_ROLE_NAME}" 
   do_append_tfvar "aws_region" "${SERVICE_REGION}"
   do_append_tfvar "service_log_retention" "${SERVICE_LOG_RETENTION}"

   # Add the other bits we need
   do_append_tfvar "service_handler" "${AWS_LAMBDA_HANDLER}"
   do_append_tfvar "zip_input" "${DIR_BUILD}" 
   do_append_tfvar "zip_output" "${FILE_PATH_SERVICE_ZIP}"


   # Dump a copy?
   if [[ ${SERVICE_ARTEFACTS_DUMP,,} == "true" ]]; then

      # Check output directory exists
      if [[ ! -d ${DIR_WORK_OUT} ]]; then
         mkdir -p ${DIR_WORK_OUT}
      fi

      # Remove file if already there
      if [[ -f ${DIR_WORK_OUT}/${FILE_NAME_TFVARS} ]]; then
         rm -rf ${DIR_WORK_OUT}/${FILE_NAME_TFVARS}
      fi

      # Copy the file
      cp ${FILE_PATH_TFVARS} ${DIR_WORK_OUT}
   fi
}


###
#
# Function to run Terraform init
#
###
do_perform_tfinit() {

   # Log start
   log_it 1 "Performing terraform init..."

   # Run init
   cd ${DIR_TERRAFORM}
   if [[ ${VERBOSE} == "TRUE" ]]; then
      log_it 2 ""
      terraform init
      RESULT=${?}
      log_it 2 ""
   else
      terraform init > /dev/null
      RESULT=${?}
   fi
   if [[ ${RESULT} -ne 0 ]]; then
      log_it 2 "*** FAILED *** ERROR reported by terraform init"
      log_it 2 ""
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
   log_it 1 "Performing terraform plan..."

   # Run plan
   cd ${DIR_TERRAFORM}
   if [[ ${VERBOSE} == "TRUE" ]]; then
      log_it 2 ""
      terraform plan -input=false -out ${FILE_PATH_TF_PLAN_OUT} -var-file ${FILE_PATH_TFVARS}
      RESULT=${?}
      log_it 2 ""
   else
      terraform plan -input=false -out ${FILE_PATH_TF_PLAN_OUT} -var-file ${FILE_PATH_TFVARS} > /dev/null
      RESULT=${?}
   fi
   if [[ ${RESULT} -ne 0 ]]; then
      log_it 2 "*** FAILED *** ERROR reported by terraform plan"
      log_it 2 ""
      exit 1;
   fi
}



#################################################
#
# Main entrypoint of script
#
#################################################
export MODE="POST"
export VERBOSE="FALSE"


# Check that we received valid options
options=':vh'
while getopts $options option
do
   case ${option} in
      h  ) do_usage; exit 0 ;;
      v  ) VERBOSE="TRUE";;
      \? ) echo ""; echo "ERROR: Unknown option: $OPTARG" 1>&2; exit 1;;
   esac
done


# Start processing
log_it 0 ""
log_it 0 "###################################################"
log_it 0 "Start of processing..."

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

   # Zip the binary
   #do_zip

   # Create TFVARS file
   do_create_tfvars

   # Perform terraform init
   do_perform_tfinit

   # Perform terraform plan
   do_perform_tfplan

fi


# Log finish
log_it 0 "End of processing..."
log_it 0 "###################################################"
