package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var sess *session.Session
var ssmsvc *ssm.SSM

const REGION string = "eu-central-1"

func createAWSClient(profile *string, verbose bool) {
	config := aws.Config{Region: aws.String(REGION), MaxRetries: aws.Int(15)}

	if verbose == true {
		config.WithLogLevel(aws.LogDebugWithRequestRetries)
	}

	sess, err := session.NewSessionWithOptions(
		session.Options{
			Config:  config,
			Profile: *profile,
		})

	if err != nil {
		fmt.Println("Session not created: ", err)
		os.Exit(127)
	}

	ssmsvc = ssm.New(sess, &config)
}

func getSecret(path *string) string {
	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(*path),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		panic(err)
	}

	value := *param.Parameter.Value
	return value
}
