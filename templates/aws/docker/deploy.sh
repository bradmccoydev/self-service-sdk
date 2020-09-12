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
export DIR_WORK="${DIR_BASE}/work"
export DIR_ZIP="${DIR_BASE}/zip"

###
# File Variables
###
export FILE_USER_INPUTS="${DIR_WORK}/inputs.sh"
export FILE_SERVICE_BINARY="${DIR_BUILD}/main"
export FILE_SERVICE_ZIP="${DIR_ZIP}/main.zip"

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
   log_it 1 "Importing user variables"

   # Read the file in
   . "${FILE_USER_INPUTS}"

   # Iterate GO_PKG_IMPORTS array & rebuild
   # to workaround weird array bug...
   GO_LIBS=()
   for i in ${!GO_PKG_IMPORTS[@]};
   do
      #echo "Processing import: ${GO_PKG_IMPORTS[$i]}"
      GO_LIBS+=(${GO_PKG_IMPORTS[$i]})
   done
}


###
#
# Function to check if variable is set
#
###
do_variable_check() {

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
# Function to perform basic sanity checks
#
###
do_sanity_checks() {

   # Log start
   log_it 1 "Performing sanity checks"

   # Check we have the user variables we need
   do_variable_check SERVICE_NAME
   do_variable_check SERVICE_DESC
   do_variable_check SERVICE_TIMEOUT
   do_variable_check SERVICE_MEMORY
   do_variable_check SERVICE_STORAGE
   do_variable_check SERVICE_RUNTIME
   do_variable_check SERVICE_ROLE_ACTION
   do_variable_check SERVICE_ROLE_NAME

   # Check source directory
   # Still to do!

}


###
#
# Function to perform build of the binary
#
###
do_build() {

   # Log start
   log_it 1 "Building microservice"
   if [[ -d ${DIR_BUILD} ]];
   then
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
   fi
 
   # Download dependencies
   log_it 2 "Downloading dependencies"
   cd ${DIR_SOURCE}
   go get ./...
   if [[ ${?} -ne 0 ]]; then
      log_it 2 "*** FAILED *** ERROR reported by go get"
      log_it 2 ""
      exit 1;
   fi
   
   # Build the binary
   log_it 2 "Performing build"
   go build -o ${FILE_SERVICE_BINARY} 
   if [[ ${?} -ne 0 ]]; then
      log_it 2 "*** FAILED *** ERROR reported by go build"
      log_it 2 ""
      exit 1;
   fi

}


###
#
# Function to perform zip of binary
#
###
do_zip() {

   # Log start
   log_it 1 "Creating zip"
   if [[ -d ${DIR_ZIP} ]];
   then
      rm -rf ${DIR_ZIP}
   fi
   mkdir -p ${DIR_ZIP}

   # Zip the binary
   zip -qr9 ${FILE_SERVICE_ZIP} ${DIR_BUILD}/*

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

   # Build the binary
   do_build

   # Zip the binary
   do_zip

   # Do Terraform stuff!
   # Still to do!

fi


# Log finish
log_it 0 "End of processing..."
log_it 0 "###################################################"
