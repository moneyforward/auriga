/*
 * Copyright 2022 Money Forward, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package slack

import (
	"context"

	"github.com/moneyforward/auriga/app/pkg/errors"
	"github.com/slack-go/slack"
)

type Client interface {
	PostMessage(ctx context.Context, channelID, message, ts string) error
	PostEphemeral(ctx context.Context, channelID, userID, ts, message string) error
	GetConversationReplies(ctx context.Context, channelID, ts string) ([]slack.Message, error)
	GetUsersInfo(ctx context.Context, userID ...string) (*[]slack.User, error)

	GetClient() *slack.Client
	GetAppUserID() string
}

type client struct {
	*slack.Client
	appUserID string
}

type Option = slack.Option

func AppTokenOption(appToken string) Option {
	return slack.OptionAppLevelToken(appToken)
}

// NewClient builds a slack client
func NewClient(botToken string, options ...Option) (Client, error) {
	c := slack.New(botToken, options...)
	at, err := c.AuthTest()
	if err != nil {
		return nil, errors.Wrap(err, "failed to authenticate test")
	}
	return &client{
		Client:    c,
		appUserID: at.UserID,
	}, nil
}

func (c *client) PostMessage(ctx context.Context, channelID, message, ts string) error {
	_, _, err := c.PostMessageContext(
		ctx,
		channelID,
		slack.MsgOptionTS(ts),
		slack.MsgOptionText(message, false),
	)

	if err != nil {
		return errors.Wrap(err, "failed to post message")
	}

	return nil
}

func (c *client) PostEphemeral(ctx context.Context, channelID, userID, ts, message string) error {
	_, err := c.PostEphemeralContext(
		ctx,
		channelID,
		userID,
		slack.MsgOptionTS(ts),
		slack.MsgOptionText(message, false),
	)

	if err != nil {
		return errors.Wrap(err, "failed to post message")
	}

	return nil
}

func (c *client) GetConversationReplies(ctx context.Context, channelID, ts string) ([]slack.Message, error) {
	params := &slack.GetConversationRepliesParameters{
		ChannelID: channelID,
		Timestamp: ts,
	}
	msgs, _, _, err := c.GetConversationRepliesContext(ctx, params)
	if err != nil {
		if err.Error() == ErrThreadNotFound.Error() {
			return nil, ErrThreadNotFound
		} else {
			return nil, errors.Wrap(err, "failed to get conversation replies")
		}
	}

	return msgs, nil
}

func (c *client) GetUsersInfo(ctx context.Context, userID ...string) (*[]slack.User, error) {
	users, err := c.GetUsersInfoContext(ctx, userID...)
	if err != nil {
		if err.Error() == ErrUserNotFound.Error() {
			return nil, ErrUserNotFound
		} else {
			return nil, errors.Wrap(err, "failed to get user info")
		}
	}

	return users, nil
}

func (c *client) GetClient() *slack.Client {
	return c.Client
}

func (c *client) GetAppUserID() string {
	return c.appUserID
}
