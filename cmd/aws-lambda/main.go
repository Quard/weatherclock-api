package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/Quard/weatherclock-api/internal/apigateway"
)

func main() {
	lambda.Start(apigateway.HandleRequest)
}
