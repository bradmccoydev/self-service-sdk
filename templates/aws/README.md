# AWS Template

This directory contains the template for building & deploying a GoLang microservice to AWS.


## Requirements

The table below lists the requirements to use the AWS template.

| Item | Requirement |
| :---: | :--- |
| 1 | [Docker](https://www.docker.com/products/docker-desktop) |
| 2 | An AWS account with appropriate permissions. |


## Template Directories

The template has the following structure:

| Directory | Purpose |
|:--:|:--|
| | The top level directory contains the docker-compose file. |
| docker | This directory contains the dockerfile and related build scripts. |
| docker/scripts | This directory contains the bash scripts used by the container. |
| docker/terraform | This directory contains the terraform scripts. used to deploy the service to AWS. |
| example | This directory contains an example microservice using the SDK. |
| work | This directory contains the files that need to be modified for your microservice such as the [inputs script](#input-variables). |


## Steps To Follow

To use this template you need to:
1. Download the template
2. Develop your microservice based on the example code.
3. Perform any local testing.
4. Update the various [input variables](#input-variables) in the {TEMPLATE_HOME}/work/inputs.sh with appropriate values
5. Update the [docker options](#docker-options) in the {TEMPLATE_HOME}/docker-compose.yml file with appropriate values
5. Start a terminal session & navigate to the directory {TEMPLATE_HOME}
6. Run docker-compose up. This will start [processing](#template-process).
5. Test the lambda microservice


## Input Variables

The template makes use of a script inputs.sh located in the work directory. This script has a number of variables that need to be populated with appropriate values to allow publishing of the lambda to AWS. The table below lists the variables, their purpose etc.

| Variable | Purpose |
|:--:|:--|
| SERVICE_NAME | The name of the service |
| SERVICE_DESC | A description for the service |
| SERVICE_TIMEOUT | The timeout value for the service |
| SERVICE_MEMORY | The memory limit for the service |
| SERVICE_STORAGE | The S3 bucket path where the service zip should be saved to |
| SERVICE_RUNTIME | The AWS Lambda runtime environment to use |
| SERVICE_ROLE_ACTION | A flag to indicate whether to create a new role or use an existing: CREATE / USE |


## Docker Options

To build & deploy your service, the container needs access to:
* your source code
* AWS credentials


### Source Code

The source code gets exposed to the container via a volume mount. Modify the directory path to point to your source code.
In the example below the container will mount the directory "example" that's relative to the docker-compose.yml directory.

*** NOTE ***
If you modify the container mount point (/tmp/source) then the deployment script will also need to be modified as it expects the source to be available there.

```
#####
#
# Source Code Mount
#
# This volume mount is used to expose the source code for the service.
# Change the mount point to the source folder on your pc.
#
#####
      - ./example:/tmp/source
      
```


### AWS Credentials

You can expose AWS credentials to the container by either:
* using a volume mount
* using environment variables

To use a volume mount, ensure your AWS credentials file is present in your home directory & uncomment the "Option 1" volume mount block.

```
#####
# Option1: Mount the ~/.aws/credentials file 
#####
#      - $HOME/.aws/credentials:/root/.aws/credentials:ro

```

If you want to use environment variables, uncomment the "Option 2" block and ensure the local variables names (ie the ones on the right side of the =) are correct.

```
#####
# Option2: Expose the credentials via environment variables
#####
#    environment:
#      - AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
#      - AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
#      - AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION}
```


### Option Details

The table below details the options that can be configured in the docker-compose.yml file.

| Option | Purpose | 
|:--:|:--|
| Work Directory Mount | To make the work directory, with the inputs script accessible to the container |
| Source Code Mount | To make the source code accessible to the container |
| AWS Credentials File Mount | To make the user's AWS credentials file accessible (in read-only mode) to the container |
| AWS Credentials Environment Variables | To make the user's AWS credentials available as environment variables |
| Run Mode | An override for the container entrypoint command. Normally this should be set to run /tmp/deploy_aws.sh. For debugging purposes you can switch to /bin/bash |


## Template Process

When you run the docker container via docker-compose up, the deploy script runs which performs the following steps:
1. Read the inputs script to get the values you provided
2. Check that all variables, source files etc are available
3. Download any dependencies (specified in the go.mod file)
4. Build the service binary 
5. Create the tfvars file to provide required values for the terraform variables
6. Run terraform init
7. Run terraform plan, using the tfvars from step 5
8. Run terraform apply, using the plan from step 7
