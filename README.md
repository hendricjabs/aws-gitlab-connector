# AWS GitLab Connector

AWS CodeBuild/CodePipeline does not nativly support GitLab Repositories. You can only build repositories from CodeCommit, GitHub and Zip-Files stored in Amazon S3. This is a sample template for a GitLab connector for AWS CodeBuild. To use it you need to set up the Lambda endpoint as described below, add it as webhook for push-events to your GitLab repository and change your CodeBuild or CodePipeline configuration so that the destination S3-Bucket is your source of code.
```bash
.
├── Makefile                    <-- Make to automate build
├── README.md                   <-- This instructions file
├── gitlab-connector            <-- Source code for a lambda function
│   ├── main.go                 <-- Lambda function code
│   ├── main_test.go            <-- Test for main function
│   ├── GitCloneAndZip.go       <-- Routine for cloning and zipping repository 
│   ├── GitCloneAndZip_test.go  <-- Test for cloning and zipping routine
|   ├── Zipper.go               <-- Routine for compressing bytes into a Zip-archive 
|   └── Zipper_test.go          <-- Test for compressing routine
├── template.yaml               <-- SAM Template
└── .gitignore                  <-- Files to ignore by git
```

## Status
[![Build Status](https://travis-ci.com/hendricjabs/aws-gitlab-connector.png)](https://travis-ci.org/hendricjabs/aws-gitlab-connector)

## Requirements

* AWS CLI already configured with Administrator permission
* [Golang](https://golang.org)
* [Makefile](https://wiki.ubuntuusers.de/Makefile/)

## Architecture
![Architecture](assets/architecture.png)

## Setup process

### Installing dependencies

You can use make to install all go dependencies or run the GO command directly.
```bash
make deps
// or
go get -u ./gitlab-connector/...
```

### Building

Golang is a staticly compiled language, meaning that in order to run it you have to build the executeable target.

You can issue the following command in a shell to build it:

```bash
make build
// or 
GOOS=linux GOARCH=amd64 go build -o gitlab-connector/gitlab-connector ./gitlab-connector
```

**NOTE**: If you're not building the function on a Linux machine, you will need to specify the `GOOS` and `GOARCH` environment variables, this allows Golang to build your function for another system architecture and ensure compatability.

### Local development

**Invoking function locally through local API Gateway**

```bash
sam local start-api
```

If the previous command ran successfully you should now be able to hit the following local endpoint to invoke your function `http://localhost:3000/gitlab`

**SAM CLI** is used to emulate both Lambda and API Gateway locally and uses the `template.yaml` to understand how to bootstrap this environment (runtime, where the source code is, etc.) - The following excerpt is what the CLI will read in order to initialize an API and its routes:

```yaml
...
Events:
    GitLab:
        Type: Api # More info about API Event Source: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#api
        Properties:
            Path: /gitlab
            Method: post
```

## Packaging and deployment

To deploy the connector you can simply use the Docker deployment.

First you need to build the image with some parameters you should define previously.
Replace 
* `<Target Bucket>` with the S3 bucket where the CloudFormation template shall be stored.
* `<Lambda Bucket>` with the S3 bucket where the Lambda deployment package shall be stored.
* `<AWS Access Key>` with your AWS access key (e.g. AKIAXXXXXXXXXXXXXXXQ)
* `<AWS Secret Access Key>` with your AWS secret access key (e.g. NMAXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXR)
* `<Region>` with the region you want to deploy the stack in

```sh
TARGET_BUCKET=<Target Bucket>
LAMBDA_BUCKET=<Lambda Bucket>
AWS_ACCESS_KEY=<AWS Access Key>
AWS_SECRET_ACCESS_KEY=<AWS Secret Access Key>
REGION=<Region>

docker build \
    --build-arg target_bucket=$TARGET_BUCKET \
    --build-arg lambda_bucket=$LAMBDA_BUCKET \
    --build-arg aws_access_key=$AWS_ACCESS_KEY \
    --build-arg aws_secret_key=$AWS_SECRET_ACCESS_KEY \ 
    --build-arg region=$REGION \
    -t hendricjabs/aws-gitlab-connector .
```

Now you can run the deployment process:
```sh
docker run hendricjabs/aws-gitlab-connector
```

After deployment is complete you can run the following command to retrieve the API Gateway Endpoint URL:

```bash
aws cloudformation describe-stacks \
    --stack-name gitlab-connector \
    --query 'Stacks[].Outputs'
``` 

### Testing

We use `testing` package that is built-in in Golang and you can simply run the following command to run our tests:

```shell
make test
// or
go test -v ./gitlab-connector/
```
# Appendix

### Golang installation

Please ensure Go 1.x (where 'x' is the latest version) is installed as per the instructions on the official golang website: https://golang.org/doc/install


## Bringing to the next level

Here are a few ideas that you can use to get more acquainted as to how this overall process works:

* ~~Create an additional API resource for Atlassian BitBucket~~ *Now supported by AWS CodeBuild.*

Next, you can use the following resources to know more about beyond hello world samples and how others structure their Serverless applications:

* [AWS Serverless Application Repository](https://aws.amazon.com/serverless/serverlessrepo/)

## Contribute
I'm glad for every feedback, improvement, recommendation, pull-request, bug report and advise.
