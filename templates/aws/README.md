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
| | The top level directory contains the docker-compose file |
| docker| This directory contains the dockerfile and docker entry script. You should not normally need to modify anything in here.|
| work| This directory contains the files that need to be modified for your microservice such as the [inputs script](#input-variables) |


## Steps To Follow

To use this template you need to:
1. Develop your microservice & confirm it is working
2. Update the various [input variables](#input-variables) in the {TEMPLATE_HOME}/work/inputs.sh with appropriate values
3. Start a terminal session & navigate to the directory {TEMPLATE_HOME}
4. Run docker-compose up. This will start [processing](#template-process).
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


## Template Process

When you run the docker container via docker-compose up, it will perform the following steps:
1. Read the inputs script to get the values you provided
2. Check that all variables, source files etc are available
3. Build the microservice binary 
4. Zip the binary
5. Upload the zip file to the development deployments S3 bucket
6. Create (or update if it exists) your lambda function
7. Insert (or update) the entry in the service table in Dynamo DB
