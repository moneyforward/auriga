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
	"fmt"

	"github.com/moneyforward/auriga/app/internal/domain/repository"
	"github.com/moneyforward/auriga/app/internal/model"
	"github.com/slack-go/slack/slackevents"
)

type SlackResponseService interface {
	ReplyEmailList(ctx context.Context, event *slackevents.AppMentionEvent, emails []*model.SlackUserEmail) error
	ReplyError(ctx context.Context, event *slackevents.AppMentionEvent, err error) error
	ReplyHelp(ctx context.Context, event *slackevents.AppMentionEvent) error
}

type slackResponseService struct {
	slackRepository repository.SlackRepository
	errorRepository repository.ErrorRepository
}

func NewSlackResponseService(factory repository.Factory) *slackResponseService {
	return &slackResponseService{
		slackRepository: factory.SlackRepository(),
		errorRepository: factory.ErrorRepository(),
	}
}

func (s *slackResponseService) ReplyEmailList(ctx context.Context, event *slackevents.AppMentionEvent, emails []*model.SlackUserEmail) error {
	msg := "参加者一覧\n"
	for _, email := range emails {
		msg += email.Email
		msg += "\n"
	}
	err := s.slackRepository.PostMessage(
		ctx,
		event.Channel,
		fmt.Sprint(msg),
		event.ThreadTimeStamp,
	)
	return err
}

func (s *slackResponseService) ReplyError(ctx context.Context, event *slackevents.AppMentionEvent, err error) error {
	var msg string
	if s.errorRepository.ErrThreadNotFound(err) {
		msg += "スレッドで呼び出してね:neko_namida:"
		return s.slackRepository.PostEphemeral(
			ctx, event.Channel, msg, event.ThreadTimeStamp, event.User,
		)
	}
	if s.errorRepository.ErrUserNotFound(err) {
		msg += "参加者はいないようです:neko_namida:"
		return s.slackRepository.PostMessage(
			ctx, event.Channel, msg, event.ThreadTimeStamp,
		)
	}
	return err
}

func (s *slackResponseService) ReplyHelp(ctx context.Context, event *slackevents.AppMentionEvent) error {
	msg := "[使い方]\n" +
		"1. スレッドで `@Auriga :sanka:` のようにAurigaを呼び出し、リアクションを指定してください。\n" +
		"2. スレッドの開始メッセージに指定のリアクションをしたユーザのメールアドレス一覧を返します。\n" +
		"3. 結果をGoogleCalenderに貼り付けると一括招待できます！"
	return s.slackRepository.PostEphemeral(
		ctx,
		event.Channel,
		fmt.Sprint(msg),
		event.ThreadTimeStamp,
		event.User,
	)
}
