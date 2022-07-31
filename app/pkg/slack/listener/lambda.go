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
	"encoding/json"
	"fmt"
	"net/http"

	pkgslack "github.com/moneyforward/auriga/app/pkg/slack"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type lambdaListener struct {
	eventHandlerFunc pkgslack.EventHandlerFunc
	signingSecretKey string
}

type handleEventRequest func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func NewLambdaListener(eventHandlerFunc pkgslack.EventHandlerFunc, signingSecretKey string) *lambdaListener {
	l := &lambdaListener{
		eventHandlerFunc: eventHandlerFunc,
		signingSecretKey: signingSecretKey,
	}

	return l
}

func (l *lambdaListener) Listen(ctx context.Context) {
	lambda.Start(l.newHandleEventRequest(ctx))
}

func (l *lambdaListener) newHandleEventRequest(ctx context.Context) handleEventRequest {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		if res, err := l.verify(request, l.signingSecretKey); err != nil {
			fmt.Println("verification failed")
			return res, err
		}

		event, err := slackevents.ParseEvent(json.RawMessage(request.Body), slackevents.OptionNoVerifyToken())
		if err != nil {
			fmt.Println("parse failed")
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}

		switch event.Data.(type) {
		case *slackevents.EventsAPICallbackEvent:
			l.eventHandlerFunc(ctx, event.InnerEvent)
		}

		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
	}
}

// verify returns the result of slack signing secret verification.
func (l *lambdaListener) verify(request events.APIGatewayProxyRequest, sc string) (events.APIGatewayProxyResponse, error) {
	body := request.Body
	header := http.Header{}
	for k, v := range request.Headers {
		header.Set(k, v)
	}

	sv, err := slack.NewSecretsVerifier(header, sc)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	if _, err := sv.Write([]byte(body)); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}
	if err := sv.Ensure(); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
