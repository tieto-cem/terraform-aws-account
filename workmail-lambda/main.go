package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/workmail"
	"fmt"
)

type LambdaEvent struct {
	Action            string   `json:"action"`
	Region            string   `json:"region"`
	OrganizationAlias string   `json:"organizationAlias"`
	GroupEmail        string   `json:"groupEmail"`
	GroupName         string   `json:"groupName"`
	UserEmails        []string `json:"userEmails"`
}

type ReturnEvent struct {
	Message string
	Warning error
}

func HandleRequest(ctx context.Context, event LambdaEvent) (*ReturnEvent, error) {

	sess := session.Must(session.NewSession())

	// Try first with Environment variables and secondly with IAM role
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
		})

	config := &aws.Config{
		Region:      &event.Region,
		Credentials: creds,
	}

	client := workmail.New(sess, config)

	createGroupAction := &ActionCreateGroup{
		WorkMailClient: client,
		event:          &event,
	}

	action, err := NewActionFactory().
		AddAction(createGroupAction).
		GetAction(&event.Action)
	if err != nil {
		return nil, err
	}

	warn, err := action.Do()
	if err != nil {
		return nil, err
	}

	return &ReturnEvent{
		Message: fmt.Sprintf("Group %s ready for usage", event.GroupName),
		Warning: warn,
	}, err
}

func main() {

	lambda.Start(HandleRequest)

	/*
	event := LambdaEvent{
		Action:            "create-group",
		Region:            "eu-west-1",
		OrganizationAlias: "cxcloud",
		GroupEmail:        "new-group@cxcloud.awsapps.com",
		GroupName:         "NewGroup",
		UserEmails: []string{
			"aws-accounts@cxcloud.awsapps.com",
		},
	}

	response, err := HandleRequest(event)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response)
	}
	*/
}
