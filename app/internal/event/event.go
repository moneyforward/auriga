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

package event

import (
	"context"

	"github.com/moneyforward/auriga/app/pkg/slack"
	"github.com/slack-go/slack/slackevents"
)

type eventHandlerFactory struct {
	appUserID      string
	handlerFactory slack.HandlerFactory
}

func NewEventHandlerFactory(appUserID string, handlerFactory slack.HandlerFactory) *eventHandlerFactory {
	return &eventHandlerFactory{
		appUserID:      appUserID,
		handlerFactory: handlerFactory,
	}
}

func (f *eventHandlerFactory) GetFunc() slack.EventHandlerFunc {
	return func(ctx context.Context, event slackevents.EventsAPIInnerEvent) {
		switch innerEv := event.Data.(type) {
		case *slackevents.AppMentionEvent:
			if innerEv.User != f.appUserID {
				f.handlerFactory.MentionEventHandler()(context.Background(), innerEv)
			}
		}
	}
}
