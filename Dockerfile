FROM golang:latest

# Set environments
ARG target_bucket
ARG lambda_bucket
ARG aws_access_key
ARG aws_secret_key
ARG username=""
ARG password=""
ARG apiKey=""
ARG region=""

ENV PATH=/root/.local/bin:$PATH
ENV REGION=${region} USERNAME=${username} PASSWORD=${password} APIKEY=${apiKey} AWS_ACCESS_KEY=$aws_access_key AWS_SECRET_KEY=${aws_secret_key} TARGET_BUCKET=${target_bucket} LAMBDA_BUCKET=${lambda_bucket}

# Install prerequisites
RUN rm /bin/sh && ln -s /bin/bash /bin/sh
RUN apt-get update
RUN apt-get install -y curl make git python-pip

ADD . src/github.com/hendricjabs/aws-gitlab-connector

# Install aws cli
RUN pip install --upgrade pip
RUN pip install awscli --upgrade

RUN source ~/.profile
RUN mkdir -p ~/.aws
RUN printf "[default]\naws_access_key_id = $AWS_ACCESS_KEY\naws_secret_access_key = $AWS_SECRET_KEY" > ~/.aws/credentials

# Build lambda binary and deploy template
RUN cd /go/src/github.com/hendricjabs/aws-gitlab-connector; make deps clean build
CMD cd /go/src/github.com/hendricjabs/aws-gitlab-connector; make S3_TARGET_BUCKET=$TARGET_BUCKET S3_LAMBDA_BUCKET=$LAMBDA_BUCKET USERNAME=$USERNAME PASSWORD=$PASSWORD API_KEY=$APIKEY REGION=$REGION deploy