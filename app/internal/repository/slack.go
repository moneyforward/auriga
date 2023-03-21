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

package repository

import (
	"context"

	slack2 "github.com/slack-go/slack"

	"github.com/moneyforward/auriga/app/pkg/slack"

	"github.com/moneyforward/auriga/app/internal/model"

	"github.com/moneyforward/auriga/app/pkg/errors"
)

type slackRepository struct {
	client slack.Client
}

func newSlackRepository(client slack.Client) *slackRepository {
	return &slackRepository{
		client: client,
	}
}

func (r *slackRepository) PostMessage(ctx context.Context, channelID, message, ts string) error {
	return r.client.PostMessage(ctx, channelID, message, ts)
}

func (r *slackRepository) PostEphemeral(ctx context.Context, channelID, message, ts, userID string) error {
	return r.client.PostEphemeral(ctx, channelID, userID, ts, message)
}

// GetParentMessage get the first message that started the thread
func (r *slackRepository) GetParentMessage(ctx context.Context, channelID, ts string) (*model.SlackMessage, error) {
	msgs, err := r.client.GetConversationReplies(ctx, channelID, ts)
	if err != nil {
		if errors.Is(err, slack.ErrThreadNotFound) {
			return nil, errThreadNotfound
		} else {
			return nil, err
		}
	}
	if len(msgs) <= 0 {
		return nil, errors.New("number of messages is zero")
	}
	parentMessage := msgs[0]
	if r.isIncompleteReaction(parentMessage.Reactions) {
		// get full reactions
		parentMessage.Reactions, err = r.client.GetReaction(ctx, channelID, ts, true)
		if err != nil {
			return nil, err
		}
	}
	var reactions []*model.SlackReaction
	for _, reaction := range parentMessage.Reactions {
		reactions = append(reactions, &model.SlackReaction{
			Name:    reaction.Name,
			UserIDs: reaction.Users,
			Count:   reaction.Count,
		})
	}
	return &model.SlackMessage{
		ChannelID: parentMessage.Channel,
		Reactions: reactions,
	}, nil
}

// isIncompleteReaction returns true if more fetches is required
// reactions[*].Count may be greater than len(reactions[*].Users), at which point a fetch is required.
func (r *slackRepository) isIncompleteReaction(reactions []slack2.ItemReaction) bool {
	for _, reaction := range reactions {
		if reaction.Count > len(reaction.Users) {
			return true
		}
	}
	return false
}

func (r *slackRepository) ListUsersEmail(ctx context.Context, userID []string) ([]*model.SlackUserEmail, error) {
	users, err := r.client.GetUsersInfo(ctx, userID...)
	if err != nil {
		if errors.Is(err, slack.ErrUserNotFound) {
			return nil, errUserNotFound
		} else {
			return nil, err
		}
	}

	var slackUsers []*model.SlackUserEmail
	for _, user := range *users {
		slackUsers = append(slackUsers, &model.SlackUserEmail{
			ID:    user.ID,
			Email: user.Profile.Email,
		})
	}
	return slackUsers, nil
}
