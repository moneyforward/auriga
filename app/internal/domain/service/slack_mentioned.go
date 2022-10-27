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
	"strings"

	"github.com/moneyforward/auriga/app/pkg/slack"

	"github.com/moneyforward/auriga/app/internal/model"
)

const (
	CommandHelp = "help"
)

type SlackMentionedService interface {
	Parse(message string) *model.MentionParseResult
}

type slackMentionedService struct {
}

func NewSlackMentionedService() *slackMentionedService {
	return &slackMentionedService{}
}

func (s *slackMentionedService) Parse(message string) *model.MentionParseResult {
	tmp := strings.Split(message, " ")
	if len(tmp) < 2 {
		// no arguments
		return &model.MentionParseResult{
			Message: message,
		}
	}
	reaction := tmp[1]
	if slack.IsReaction(reaction) {
		return &model.MentionParseResult{
			Message: message,
			Reaction: slack.ExtractReactionName(
				slack.RemoveSkinToneFromReaction(reaction)),
		}
	}
	// invalid argument (not emoji)
	return &model.MentionParseResult{
		Message: message,
		Command: CommandHelp,
	}
}
