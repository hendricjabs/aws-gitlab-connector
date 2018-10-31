CFLAGS=-g
export CFLAGS

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
	aws cloudformation package --template-file ./template.yaml  --s3-bucket $(S3_TEMPLATE_BUCKET)  --output-template-file packaged-template.yaml
	aws cloudformation deploy --template-file ./packaged-template.yaml --stack-name $(STACK_NAME) --parameter-overrides BucketName=$S3_TARGET_BUCKET --capabilities CAPABILITY_IAM --region eu-central-1
