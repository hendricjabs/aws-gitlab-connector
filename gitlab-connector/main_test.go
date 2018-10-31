package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	t.Run("Empty Body", func(t *testing.T) {
		_, err := handler(events.APIGatewayProxyRequest{Body: ""})
		if err == nil {
			t.Fatal("Error: Should fail on empty body")
		}
	})
}
