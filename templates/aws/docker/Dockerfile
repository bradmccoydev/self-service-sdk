# Use the specified Hashicorp Terraform version
FROM amazonlinux:2.0.20200722.0


# Install required os components
RUN yum install -y gzip tar wget unzip zip


# Install GoLang
ARG GO_VERS=1.15.1
ARG GO_URL="https://dl.google.com/go/go${GO_VERS}.linux-amd64.tar.gz"
RUN wget -O /tmp/go.tar.gz --quiet ${GO_URL} \
 && tar -C /usr/local -xzf /tmp/go.tar.gz 


# Install Terraform
ARG TF_VERS=0.13.2
ARG TF_URL="https://releases.hashicorp.com/terraform/${TF_VERS}/terraform_${TF_VERS}_linux_amd64.zip"
RUN wget -O /tmp/tf.zip --quiet ${TF_URL} \
 && unzip /tmp/tf.zip -d /usr/local/bin


# Environment variables for the container
ENV PATH ${PATH}:/usr/local/go/bin


# Copy the required scripts
COPY ./docker/docker-entrypoint.sh "/docker-entrypoint.sh"
COPY ./docker/deploy.sh "/tmp/deploy.sh"


# Set required permissions
RUN chmod +x "/docker-entrypoint.sh" \
 && chmod +x "/tmp/deploy.sh"


# Set the work directory
WORKDIR /tmp/


# Set the entrypoint
ENTRYPOINT ["/docker-entrypoint.sh"]
