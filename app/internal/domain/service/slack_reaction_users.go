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

package service

import (
	"context"

	"github.com/moneyforward/auriga/app/pkg/slack"

	"github.com/moneyforward/auriga/app/pkg/slice"

	repository2 "github.com/moneyforward/auriga/app/internal/domain/repository"
	"github.com/moneyforward/auriga/app/internal/model"
)

type SlackReactionUsersService interface {
	// ListUsersEmailByReaction get the email address of the users
	// who reacted to the parent message associated with the thread
	ListUsersEmailByReaction(ctx context.Context, channelID, ts, reactionName string) ([]*model.SlackUserEmail, error)
}

type slackReactionUsersService struct {
	slackRepository repository2.SlackRepository
}

func NewSlackReactionUsersService(factory repository2.Factory) *slackReactionUsersService {
	return &slackReactionUsersService{
		slackRepository: factory.SlackRepository(),
	}
}

func (s *slackReactionUsersService) ListUsersEmailByReaction(ctx context.Context, channelID, ts, reactionName string) ([]*model.SlackUserEmail, error) {
	msg, err := s.slackRepository.GetParentMessage(ctx, channelID, ts)
	if err != nil {
		return nil, err
	}
	inviteUserIDs := s.getReactionUserIDs(ctx, msg.Reactions, reactionName)
	inviteUserEmails, err := s.slackRepository.ListUsersEmail(ctx, inviteUserIDs)
	if err != nil {
		return nil, err
	}
	return inviteUserEmails, nil
}

// getReactionUserIDs get reaction users by reactionName
func (s *slackReactionUsersService) getReactionUserIDs(ctx context.Context, reactions []*model.SlackReaction, reactionName string) []string {
	var userIDs []string
	var targetReactions []*model.SlackReaction
	for _, reaction := range reactions {
		rn := slack.ExtractReactionName(reaction.Name)
		if slack.RemoveSkinToneFromReaction(rn) == reactionName {
			targetReactions = append(targetReactions, reaction)
		}
	}
	if len(targetReactions) == 0 {
		return userIDs // no reaction members
	}
	for _, tr := range targetReactions {
		userIDs = append(userIDs, tr.UserIDs...)
	}
	return slice.ToStringSet(userIDs)
}
