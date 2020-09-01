# Self Service SDK

This repository contains the SDK for Self Service. It is provided as a GoLang module.


## Requirements

The table below lists the requirements to use the SDK.

| Item | Requirement |
| :---: | :--- |
| 1 | Go >= v1.14 |
| 2 | The test framework expects certain environment variables to be available. See [below](#testing). |


## Testing

Certain environment variables need to be set for the various unit tests to run. These are summarised in the table below.

| Variable | Notes |
| :---: | :--- |
| TESTING_AWS_ACCESS_KEY_ID | A valid AWS Access Key that can be used to connect to AWS |
| TESTING_AWS_SECRET_ACCESS_KEY | The secret for the above AWS access key |
| TESTING_AWS_DEFAULT_REGION | The default AWS region to use |
