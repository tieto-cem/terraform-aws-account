package main

import (
	"github.com/aws/aws-sdk-go/service/workmail"
	"github.com/aws/aws-sdk-go/aws"
)

// WorkMail
type WorkMail struct {
	client *workmail.WorkMail
	event  *LambdaEvent
}

func (wm *WorkMail) GetOrganizationID() (*string, error) {
	var organizationID *string

	// Find organization ID
	input := workmail.ListOrganizationsInput{
		MaxResults: aws.Int64(100),
	}
	organizations, err := wm.client.ListOrganizations(&input)
	if err != nil {
		return organizationID, err
	}

	for _, organization := range organizations.OrganizationSummaries {
		if *organization.Alias == wm.event.OrganizationAlias {
			organizationID = organization.OrganizationId
		}
	}

	return organizationID, err
}

func (wm *WorkMail) getUsers(organizationID *string) ([]*workmail.User, error) {
	listUserInput := workmail.ListUsersInput{
		OrganizationId: organizationID,
		MaxResults:     aws.Int64(100),
	}
	userList, err := wm.client.ListUsers(&listUserInput)

	return userList.Users, err
}

func (wm *WorkMail) getGroups(organizationID *string) (*workmail.ListGroupsOutput, error) {
	groupsInput := workmail.ListGroupsInput{
		MaxResults:     aws.Int64(100),
		OrganizationId: organizationID,
	}
	groups, err := wm.client.ListGroups(&groupsInput)
	return groups, err
}

func (wm *WorkMail) GetUserIDs(organizationID *string) ([]*workmail.User, error) {
	var users []*workmail.User

	allUsers, err := wm.getUsers(organizationID)
	if err != nil {
		return users, err
	}

	for _, user := range allUsers {
		for _, userEmail := range wm.event.UserEmails {
			if user.Email != nil && *user.Email == userEmail {
				users = append(users, user)
			}
		}
	}
	return users, err
}

func (wm *WorkMail) GetGroupIDFromEmail(organizationID *string) (*string, error) {
	var groupID *string

	groups, err := wm.getGroups(organizationID)
	if err != nil {
		return groupID, err
	}

	for _, group := range groups.Groups {

		if group.Email != nil && *group.Email == wm.event.GroupEmail {
			groupID = group.Id
		}
	}

	return groupID, err
}

func (wm *WorkMail) GetGroupFromName(organizationID *string) (*workmail.Group, error) {

	var rGroup *workmail.Group

	groups, err := wm.getGroups(organizationID)
	if err != nil {
		return nil, err
	}

	for _, group := range groups.Groups {
		if *group.Name == wm.event.GroupName {
			rGroup = group
		}
	}

	return rGroup, err
}

func (wm *WorkMail) EnableGroup(organizationID, groupID *string) error {
	rigisterGroup := workmail.RegisterToWorkMailInput{
		Email:          &wm.event.GroupEmail,
		OrganizationId: organizationID,
		EntityId:       groupID,
	}

	_, err := wm.client.RegisterToWorkMail(&rigisterGroup)
	return err
}

func (wm *WorkMail) CreateGroup(organizationID *string) (*string, error) {
	groupInput := workmail.CreateGroupInput{
		OrganizationId: organizationID,
		Name:           &wm.event.GroupName,
	}
	group, err := wm.client.CreateGroup(&groupInput)
	return group.GroupId, err
}

func (wm *WorkMail) AssociateMembersToGroup(organizationID, groupID *string, members []*workmail.User) error {
	var err error
	for _, member := range members {
		associateMemberInput := workmail.AssociateMemberToGroupInput{
			OrganizationId: organizationID,
			GroupId:        groupID,
			MemberId:       member.Id,
		}
		_, err = wm.client.AssociateMemberToGroup(&associateMemberInput)
		if err != nil {
			break
		}
	}
	return err
}
