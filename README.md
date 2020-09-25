# Self Service SDK

This repository contains the SDK for Self Service. It is provided as a GoLang module.


## Dirctory Structure

The table below explains the layout of this directory.

| Directory | Purpose |
| :--- | :--- |
| /internal | This directory contains utility functions for internal testing etc |
| /sdk | This directory contains the [SDK packages](sdk/README.md) |
| /templates | This directory contains [templates](templates.README.md) for using the SDK |


## Requirements

The table below lists the requirements to use the SDK.

| Item | Requirement |
| :---: | :--- |
| 1 | Go >= v1.15 |
| 2 | An AWS account with appropriate permissions. See [below](#aws). |
| 3 | Certain environment variables to drive the testing framework. See [below](#testing). |


## AWS

To enable testing of the AWS functionality an AWS account is required. The table below lists the permissions required.

| Object | Permission | Notes |
| :--- | :--- | :--- |
| DynamoDB | ? | |
| Secrets Manager | ? | |
| | | |


## Testing

Certain environment variables need to be set for the various unit tests to run. These are summarised in the table below.

| Variable | Notes |
| :--- | :--- |
| TESTING_AWS_ENABLED | Set to TRUE if the AWS related tests should be run |
| TESTING_AWS_ACCESS_KEY_ID | A valid AWS Access Key that can be used to connect to AWS |
| TESTING_AWS_SECRET_ACCESS_KEY | The secret for the above AWS access key |
| TESTING_AWS_DEFAULT_REGION | The default AWS region to use |
| TESTING_AWS_USER_ID | The User ID of the above AWS user. This is used for validating a successful login |
