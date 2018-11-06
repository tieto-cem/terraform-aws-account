package main

import (
	"github.com/aws/aws-sdk-go/service/workmail"
	"fmt"
)

// ActionList only list Jira issues from git commits
type ActionCreateGroup struct {
	WorkMailClient *workmail.WorkMail
	event          *LambdaEvent
}

// Name return the name of the action
func (cg *ActionCreateGroup) Name() string {
	return "create-group"
}

func (cg *ActionCreateGroup) Do() (error, error) {

	wMail := WorkMail{
		client: cg.WorkMailClient,
		event:  cg.event,
	}

	organizationID, err := wMail.GetOrganizationID()
	if err != nil {
		return nil, err
	}
	if organizationID == nil {
		return nil, fmt.Errorf("organazation %s doesn't exist", cg.event.OrganizationAlias)
	}

	users, err := wMail.GetUserIDs(organizationID)
	if err != nil {
		return nil, err
	}

	groupID, err := wMail.GetGroupIDFromEmail(organizationID)
	if err != nil {
		return nil, err
	}
	if groupID != nil {
		return fmt.Errorf("email address %s is already in use", cg.event.GroupEmail), nil
	}

	group, err := wMail.GetGroupFromName(organizationID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		// Create new group
		groupID, err = wMail.CreateGroup(organizationID)
		if err != nil {
			return nil, err
		}
		err := wMail.EnableGroup(organizationID, groupID)
		if err != nil {
			return nil, err
		}
	} else if *group.State == "DISABLED" {
		groupID = group.Id
		err := wMail.EnableGroup(organizationID, groupID)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("group %s is already in use", cg.event.GroupName)
	}

	err = wMail.AssociateMembersToGroup(organizationID, groupID, users)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
