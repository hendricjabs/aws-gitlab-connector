package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	log "github.com/sirupsen/logrus"
)

// Repository: struct for repository information
type Repository struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

// Commit: struct for commit information
type Commit struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// WebHookRequest: Location of the repository and other relevant information
type WebHookRequest struct {
	Repository Repository `json:"repository"`
	Commits    []Commit   `json:"commits"`
}

// WebHookResponse: Location of the zip file in S3
type WebHookResponse struct {
	Bucket string
	Key    string
}

var (
	API_KEY   = os.Getenv("API_KEY")
	S3_BUCKET = os.Getenv("S3_BUCKET")
	S3_REGION = os.Getenv("S3_REGION")
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout for CloudWatch
	log.SetOutput(os.Stdout)

	// Only log the info severity or above.
	log.SetLevel(log.InfoLevel)
}

// CheckIfError: check if an error occurred and log if necessary
func CheckIfError(err error) bool {
	if err != nil {
		log.WithField("Error", err).Errorf("An error occurred")
		return true
	}
	return false
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	webHookRequest := &WebHookRequest{}

	// unmarshal body to struct object
	err := json.Unmarshal([]byte(request.Body), webHookRequest)
	if CheckIfError(err) {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 501,
		}, err
	}

	log.WithFields(log.Fields{
		"URL":  *webHookRequest,
		"Body": request.Body,
	}).Info("Got following Parameter:")

	// clone repository and put into zip file
	file, err := GitCloneAndZip(&GitCloneAndZipInput{apiKey: API_KEY, repositoryUrl: webHookRequest.Repository.URL})
	if CheckIfError(err) {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 502,
		}, err
	}

	// generate filename from repository name and id of last commit
	fileName := fmt.Sprintf("%s-%s.zip", webHookRequest.Repository.Name, webHookRequest.Commits[0].ID)

	// upload the zip file to S3
	err = addFileToS3(file.File, fileName, file.Size)
	if CheckIfError(err) {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 503,
		}, err
	}
	// marshal response to json
	response, err := json.Marshal(WebHookResponse{Bucket: S3_BUCKET, Key: fileName})
	if CheckIfError(err) {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 503,
		}, err
	}

	// Return the location data to API Gateway
	return events.APIGatewayProxyResponse{
		Body:       string(response),
		StatusCode: 200,
	}, err
}

// addFileToS3: puts the given file into the defined bucket to the given path
// if an error occurs it is returned for further handling
func addFileToS3(file []byte, path string, size int64) error {
	log.WithField("Path", path).Info("Got Parameters for S3 Put")
	s, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION)})
	// create buffer
	buffer := make([]byte, size)

	// upload object to s3
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(S3_BUCKET),
		Key:                aws.String(path),
		ACL:                aws.String("private"),
		Body:               bytes.NewReader(file),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
	})
	return err
}

func main() {
	lambda.Start(handler)
}
