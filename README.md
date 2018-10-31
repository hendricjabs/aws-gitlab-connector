# git-connector

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

AWS Lambda Python runtime requires a flat folder with all dependencies including the application. SAM will use `CodeUri` property to know where to look up for both application and dependencies:

```yaml
...
    FirstFunction:
        Type: AWS::Serverless::Function
        Properties:
            CodeUri: gitlab-connector/
            ...
```

First and foremost, we need a S3 bucket where we can upload our Lambda functions packaged as ZIP and another S3 bucket which is used as target for the packaged Git-Repository - If you don't have the S3 buckets to store code artifacts and result files then this is a good time to create one:

```bash
aws s3 mb s3://BUCKET_NAME
```

Next you can use `make` to deploy this setup. Replace `<TEMPLATE BUCKET>` with the bucket where the Lambda function is stored, <CLOUDFORMATION STACK NAME> with a account-unique name for the CloudFormation stack and <TARGET BUCKET> with the S3 bucket name of the bucket, where the result ZIP files will be stored
```bash
make deploy S3_LAMBDA_BUCKET=<LAMBDA BUCKET> STACK_NAME=<CLOUDFORMATION STACK NAME> S3_TARGET_BUCKET=<TARGET BUCKET>
```
> **See [Serverless Application Model (SAM) HOWTO Guide](https://github.com/awslabs/serverless-application-model/blob/master/HOWTO.md) for more details in how to get started.**

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

* Create an additional API resource for Atlassian BitBucket
* Update unit test to capture that
* Package & Deploy

Next, you can use the following resources to know more about beyond hello world samples and how others structure their Serverless applications:

* [AWS Serverless Application Repository](https://aws.amazon.com/serverless/serverlessrepo/)

## Contribute
I'm glad for every feedback, improvement, recommendation, pull-request, bug report and advise.
