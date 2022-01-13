package service

import (
	"dependabot/internal/config"
	"dependabot/internal/errors"
	"fmt"

	"github.com/slack-go/slack"
)

type SlackNotificationService interface {
	NotifyMerge(channelId string, mergeRequestURL string, repositoryName string, packageName string) error
}

type slackNotificationService struct {
	config      *config.Main
	slackClient *slack.Client
}

func (s *slackNotificationService) NotifyMerge(channelId string, mergeRequestURL string, repositoryName string, packageName string) error {
	_, _, err := s.slackClient.PostMessage(channelId, mergeSubmittedMessage(mergeRequestURL, repositoryName, packageName))

	if err != nil {
		return errors.NewOperationError(err, "failure when post message to slack channel %s", channelId)
	}
	return nil
}

func mergeSubmittedMessage(mergeRequestURL string, repositoryName string, packageName string) slack.MsgOption {
	text := fmt.Sprintf(":github-reviewed: [Dependency Bump] in repository %s for %s dependencies", repositoryName, packageName)
	txtBlockObject := slack.NewTextBlockObject("mrkdwn", text, false, false)

	buttonTextObject := slack.NewTextBlockObject("plain_text", "Start Review!", true, false)
	reviewButton := slack.NewButtonBlockElement("review-merge-request", mergeRequestURL, buttonTextObject)
	reviewButton.URL = mergeRequestURL

	return slack.MsgOptionBlocks(slack.NewSectionBlock(txtBlockObject, nil, slack.NewAccessory(reviewButton)))
}

func NewSlackNotificationService(
	config *config.Main,
	slackClient *slack.Client,
) SlackNotificationService {
	return &slackNotificationService{
		config:      config,
		slackClient: slackClient,
	}
}
