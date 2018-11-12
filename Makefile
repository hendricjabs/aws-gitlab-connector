CFLAGS=-g
export CFLAGS

API_KEY = ''
USERNAME = ''
PASSWORD = ''
REGION = 'eu-central-1'

.PHONY: deps clean build

deps:
	go get -u ./gitlab-connector/...

clean:
	rm -rf ./gitlab-connector/git-connector

build:
	GOOS=linux GOARCH=amd64 go build -o gitlab-connector/gitlab-connector ./gitlab-connector

test:
	go test ./gitlab-connector

deploy:
	if ! aws s3 ls "$(S3_LAMBDA_BUCKET)" 2> /dev/null; then aws s3 mb s3://$(S3_LAMBDA_BUCKET); fi
	aws cloudformation package --template-file ./template.yaml --s3-bucket $(S3_LAMBDA_BUCKET) --output-template-file packaged-template.yaml
	aws cloudformation deploy --template-file ./packaged-template.yaml --stack-name "GitLabConnector" --parameter-overrides BucketName=$(S3_TARGET_BUCKET) Username=$(USERNAME) Password=$(PASSWORD) ApiKey=$(API_KEY) --capabilities CAPABILITY_IAM --region $(REGION)
