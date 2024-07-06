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
	"strings"

	"github.com/moneyforward/auriga/app/pkg/slice"

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

const (
	// lineSizeOfPostEmailList is limit number of line when Auriga reply Email List message.
	// According to Slack API Documentation, the limit number of characters of postMessageAPI is 4000 for the best results.
	// If the average length of email addresses sent at one time exceeds approximately 80, there is a risk of error within this method.
	// but, it is the number we can afford.
	lineSizeOfPostEmailList = 50
)

// postEmailList method posts emailList using slack postMessageAPI.
// The chunkedLines are generated and requested for each chunk,
// because of considering the limit the number of characters of slackAPI.
func (s *slackResponseService) postEmailList(ctx context.Context, channelID string, emails []*model.SlackUserEmail, ts string, userID string) error {
	lines := append(make([]string, 0, len(emails)+1), "参加者一覧")
	for _, email := range emails {
		lines = append(lines, email.Email)
	}
	chunkedLines := slice.SplitStringSliceInChunks(lines, lineSizeOfPostEmailList)
	for _, chunkedLine := range chunkedLines {
		err := s.slackRepository.PostEphemeral(ctx, channelID, strings.Join(chunkedLine, "\n"), ts, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *slackResponseService) ReplyEmailList(ctx context.Context, event *slackevents.AppMentionEvent, emails []*model.SlackUserEmail) error {
	if len(emails) <= lineSizeOfPostEmailList-1 {
		var b strings.Builder
		b.WriteString("参加者一覧")
		for _, email := range emails {
			b.WriteString("\n" + email.Email)
		}
		return s.slackRepository.PostEphemeral(ctx, event.Channel, b.String(), event.ThreadTimeStamp, event.User)
	}

	return s.postEmailList(ctx, event.Channel, emails, event.ThreadTimeStamp, event.User)
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
		return s.slackRepository.PostEphemeral(
			ctx, event.Channel, msg, event.ThreadTimeStamp, event.User,
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
