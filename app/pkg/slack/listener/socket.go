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

package listener

import (
	"context"
	"fmt"

	"github.com/moneyforward/auriga/app/pkg/slack"

	"github.com/slack-go/slack/slackevents"

	"github.com/slack-go/slack/socketmode"
)

type socketListener struct {
	socketClient     slack.SocketClient
	eventHandlerFunc slack.EventHandlerFunc
}

func NewSocketListener(socketClient slack.SocketClient, eventHandlerFunc slack.EventHandlerFunc) *socketListener {
	return &socketListener{
		socketClient:     socketClient,
		eventHandlerFunc: eventHandlerFunc,
	}
}

func (l *socketListener) Listen(ctx context.Context) {
	go l.listen(ctx)
	l.wait()
}

func (l *socketListener) listen(ctx context.Context) {
	for ev := range l.socketClient.Events() {
		switch ev.Type {
		case socketmode.EventTypeEventsAPI:
			l.socketClient.Ack(*ev.Request)
			payload := ev.Data.(slackevents.EventsAPIEvent)
			switch payload.Type {
			case slackevents.CallbackEvent:
				l.eventHandlerFunc(ctx, payload.InnerEvent)
			}
		default:
			l.socketClient.Debugf("Skipped: %v", ev.Type)
		}
	}
}

// wait for the SIGNAL
func (l *socketListener) wait() {
	if err := l.socketClient.Run(); err != nil {
		fmt.Print("warn: error in socketCli.Run")
	}
}
