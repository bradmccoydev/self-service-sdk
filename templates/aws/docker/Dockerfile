# Use the specified AWS Linux image
FROM amazonlinux:2.0.20200722.0


# Install required os utilities
RUN yum install -y gzip tar wget unzip zip


# Install GoLang
ARG GO_VERS=1.15.1
ARG GO_URL="https://dl.google.com/go/go${GO_VERS}.linux-amd64.tar.gz"
RUN wget -O /tmp/go.tar.gz --quiet ${GO_URL} \
 && tar -C /usr/local -xzf /tmp/go.tar.gz 


# Install Terraform
ARG TF_VERS=0.13.3
ARG TF_URL="https://releases.hashicorp.com/terraform/${TF_VERS}/terraform_${TF_VERS}_linux_amd64.zip"
RUN wget -O /tmp/tf.zip --quiet ${TF_URL} \
 && unzip /tmp/tf.zip -d /usr/local/bin


# Environment variables for the container
ENV PATH ${PATH}:/usr/local/go/bin


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
