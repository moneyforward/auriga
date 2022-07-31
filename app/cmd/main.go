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

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/moneyforward/auriga/app/internal/event"

	"github.com/moneyforward/auriga/app/pkg/slack/listener"

	"github.com/joho/godotenv"

	"github.com/moneyforward/auriga/app/internal/handler"

	"github.com/moneyforward/auriga/app/pkg/slack"
)

const (
	slackBotTokenKey = "SLACK_BOT_TOKEN"
	slackAppTokenKey = "SLACK_APP_TOKEN"

	slackSigningSecretKey = "SLACK_SIGNING_SECRET"
)

var (
	isDebug bool
)

func run(ctx context.Context) error {
	flag.BoolVar(&isDebug, "debug", false, "debug mode")
	flag.Parse()

	var eventListener slack.Listener

	slackClientOptions := []slack.Option{}
	if isDebug {
		if err := godotenv.Load(); err != nil {
			return fmt.Errorf("load env file failed: %v", err)
		}
		slackClientOptions = append(slackClientOptions, slack.AppTokenOption(os.Getenv(slackAppTokenKey)))
	}

	slackClient, err := slack.NewClient(os.Getenv(slackBotTokenKey), slackClientOptions...)
	if err != nil {
		return err
	}

	handlerFactory := handler.NewHandlerFactory(slackClient)
	eventHandlerFactory := event.NewEventHandlerFactory(slackClient.GetAppUserID(), handlerFactory)

	if isDebug {
		socketClient := slack.NewSocketClient(slackClient, true)
		eventListener = listener.NewSocketListener(socketClient, eventHandlerFactory.GetFunc())
	} else {
		fmt.Println("this is prod mode!!")
		eventListener = listener.NewLambdaListener(eventHandlerFactory.GetFunc(), os.Getenv(slackSigningSecretKey))
	}

	eventListener.Listen(ctx)

	return nil
}

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		panic(err)
	}
}
