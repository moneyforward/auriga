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

package handler

import (
	"context"
	"log"

	"github.com/moneyforward/auriga/app/internal/domain/service"

	"github.com/moneyforward/auriga/app/internal/repository"
	"github.com/moneyforward/auriga/app/pkg/slack"
	"github.com/slack-go/slack/slackevents"
)

type AppMentionHandler interface {
	GetFunc() slack.MentionEventHandler
}

type appMentionHandler struct {
	slackReactionUsersService service.SlackReactionUsersService
	slackResponseService      service.SlackResponseService
	slackMentionedService     service.SlackMentionedService
}

func NewAppMentionHandler(client slack.Client) *appMentionHandler {
	factory := repository.NewFactory(client)
	return &appMentionHandler{
		slackReactionUsersService: service.NewSlackReactionUsersService(factory),
		slackResponseService:      service.NewSlackResponseService(factory),
		slackMentionedService:     service.NewSlackMentionedService(),
	}
}

func (h *appMentionHandler) GetFunc() slack.MentionEventHandler {
	return func(ctx context.Context, event *slackevents.AppMentionEvent) {
		reaction := h.slackMentionedService.Parse(event.Text)
		if reaction.Command == service.CommandHelp {
			if err := h.slackResponseService.ReplyHelp(ctx, event); err != nil {
				log.Printf("Failed to reply help: %v", err)
			}
			return
		}
		if reaction.Reaction == "" {
			if err := h.slackResponseService.ReplyHelp(ctx, event); err != nil {
				log.Printf("Failed to reply help: %v", err)
			}
			return
		}
		emails, err := h.slackReactionUsersService.ListUsersEmailByReaction(ctx, event.Channel, event.ThreadTimeStamp, reaction.Reaction)
		if err != nil {
			if err = h.slackResponseService.ReplyError(ctx, event, err); err != nil {
				log.Printf("Failed to reply error: %v", err)
			}
			return
		}
		err = h.slackResponseService.ReplyEmailList(ctx, event, emails)
		if err != nil {
			log.Printf("Failed to reply: %v", err)
		}
	}
}
