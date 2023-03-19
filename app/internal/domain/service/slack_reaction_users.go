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

const (
	// ChunkSizeOfChunkedListUserEmail chunk size of calling slackRepository.ListUsersEmail
	ChunkSizeOfChunkedListUserEmail = 20
)

// chunkedListUsersEmail splits userID array into chunks,
//　and calls slackRepository.ListUsersEmail for each chunk.
func (s *slackReactionUsersService) chunkedListUsersEmail(ctx context.Context, userIDs []string) ([]*model.SlackUserEmail, error) {
	chunkedUserIDsList := slice.SplitStringSliceInChunks(userIDs, ChunkSizeOfChunkedListUserEmail)
	slackUserEmails := make([]*model.SlackUserEmail, 0, len(chunkedUserIDsList))
	for _, chunkedUserIDs := range chunkedUserIDsList {
		userEmails, err := s.slackRepository.ListUsersEmail(ctx, chunkedUserIDs)
		if err != nil {
			return nil, err
		}
		slackUserEmails = append(slackUserEmails, userEmails...)
	}
	return slackUserEmails, nil
}

func (s *slackReactionUsersService) ListUsersEmailByReaction(ctx context.Context, channelID, ts, reactionName string) ([]*model.SlackUserEmail, error) {
	msg, err := s.slackRepository.GetParentMessage(ctx, channelID, ts)
	if err != nil {
		return nil, err
	}
	if s.isReactionRefetchNeeded(msg.Reactions) {
		msg.Reactions, err = s.getFullReactions(ctx, channelID, ts)
		if err != nil {
			return nil, err
		}
	}
	reactedUserIDs := s.getReactionUserIDs(ctx, msg.Reactions, reactionName)
	reactedUserEmails, err := s.chunkedListUsersEmail(ctx, reactedUserIDs)
	if err != nil {
		return nil, err
	}
	return reactedUserEmails, nil
}

// isReactionRefetchNeeded returns true if more fetches is required
func (s *slackReactionUsersService) isReactionRefetchNeeded(reactions []*model.SlackReaction) bool {
	for _, r := range reactions {
		if r.Count > len(r.UserIDs) {
			return true
		}
	}
	return false
}

// getReactionUserIDs get reaction users by reactionName
func (s *slackReactionUsersService) getReactionUserIDs(ctx context.Context, reactions []*model.SlackReaction, reactionName string) []string {
	var userIDs []string
	var targetReactions []*model.SlackReaction
	// filtered reactions by reactionName
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

func (s *slackReactionUsersService) getFullReactions(ctx context.Context, channelID, ts string) ([]*model.SlackReaction, error) {
	return s.slackRepository.GetReactions(ctx, channelID, ts, true)
}
