# Use the specified Debian image
FROM debian:sid-slim


# Build arguments
ARG GO_VERS=1.15.1
ARG TF_VERS=0.13.3


# Install required os utilities
RUN apt-get -qq update \
 && apt-get -qq install curl jq unzip


# Install GoLang
ARG GO_URL="https://dl.google.com/go/go${GO_VERS}.linux-amd64.tar.gz"
RUN curl ${GO_URL} -o /tmp/go.tar.gz -s \
 && tar -C /usr/local -xzf /tmp/go.tar.gz 


# Update PATH to include GoLang
ENV PATH ${PATH}:/usr/local/go/bin


# Install Terraform
ARG TF_URL="https://releases.hashicorp.com/terraform/${TF_VERS}/terraform_${TF_VERS}_linux_amd64.zip"
RUN curl ${TF_URL} -o /tmp/tf.zip -s \
 && unzip /tmp/tf.zip -d /usr/local/bin


# Copy the required files
COPY ./docker/scripts/docker-entrypoint.sh "/docker-entrypoint.sh"
COPY ./docker/scripts/deploy_aws.sh "/tmp/deploy_aws.sh"
COPY ./docker/terraform /tmp/terraform


# Set required permissions
RUN chmod +x "/docker-entrypoint.sh" \
 && chmod +x "/tmp/deploy_aws.sh"


# Set the work directory
WORKDIR /tmp/


# Set the entrypoint
ENTRYPOINT ["/docker-entrypoint.sh"]
