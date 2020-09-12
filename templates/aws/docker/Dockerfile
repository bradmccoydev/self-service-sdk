# Use the specified Hashicorp Terraform version
FROM amazonlinux:2.0.20200722.0


# Install required os components
RUN yum install -y gzip tar wget zip


# Install GoLang
ARG GOLANG_VERS=1.15.1
ARG GO_URL="https://dl.google.com/go/go${GOLANG_VERS}.linux-amd64.tar.gz"
RUN wget -O /tmp/go.tar.gz ${GO_URL} \
 && tar -C /usr/local -xzf /tmp/go.tar.gz 


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
