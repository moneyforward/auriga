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
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mock_repository "github.com/moneyforward/auriga/app/internal/domain/repository/mock"

	"github.com/moneyforward/auriga/app/internal/model"
	"github.com/slack-go/slack/slackevents"
)

func Test_slackResponseService_postEmailList(t *testing.T) {
	type args struct {
		channelID               string
		message                 string
		ts                      string
		lineSizeOfPostEmailList int
	}
	tests := []struct {
		name    string
		args    args
		prepare func(msr *mock_repository.MockSlackRepository)
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				channelID:               "cid",
				message:                 "参加者一覧\nuser01@example.com\nuser02@example.com\nuser03@example.com",
				ts:                      "ts",
				lineSizeOfPostEmailList: 2,
			},
			prepare: func(msr *mock_repository.MockSlackRepository) {
				gomock.InOrder(
					msr.EXPECT().PostMessage(
						gomock.Any(), "cid", "参加者一覧\nuser01@example.com", "ts").
						Return(nil),
					msr.EXPECT().PostMessage(
						gomock.Any(), "cid", "user02@example.com\nuser03@example.com", "ts").
						Return(nil),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			msr := mock_repository.NewMockSlackRepository(ctrl)
			mer := mock_repository.NewMockErrorRepository(ctrl)
			if tt.prepare != nil {
				tt.prepare(msr)
			}
			s := &slackResponseService{
				slackRepository: msr,
				errorRepository: mer,
			}
			if err := s.postEmailList(ctx, tt.args.channelID, tt.args.message, tt.args.ts, tt.args.lineSizeOfPostEmailList); (err != nil) != tt.wantErr {
				t.Errorf("postEmailList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_slackErrorResponseService_ReplyEmailList(t *testing.T) {
	type args struct {
		event  *slackevents.AppMentionEvent
		emails []*model.SlackUserEmail
	}
	tests := []struct {
		name    string
		args    args
		prepare func(msr *mock_repository.MockSlackRepository)
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				event: &slackevents.AppMentionEvent{
					Channel:         "sampleChannel",
					ThreadTimeStamp: "sampleThreadTimeStamp",
				},
				emails: []*model.SlackUserEmail{
					{Email: "sample01@example.com"},
					{Email: "sample02@example.com"},
				},
			},
			prepare: func(msr *mock_repository.MockSlackRepository) {
				msr.EXPECT().PostMessage(gomock.Any(), "sampleChannel",
					"参加者一覧\nsample01@example.com\nsample02@example.com\n",
					"sampleThreadTimeStamp").Return(nil)
			},
		},
		{
			name: "NG: error in slackRepository/PostMessage",
			args: args{
				event: &slackevents.AppMentionEvent{
					Channel:         "sampleChannel",
					ThreadTimeStamp: "sampleThreadTimeStamp",
				},
				emails: []*model.SlackUserEmail{
					{Email: "sample01@example.com"},
					{Email: "sample02@example.com"},
				},
			},
			prepare: func(msr *mock_repository.MockSlackRepository) {
				msr.EXPECT().PostMessage(gomock.Any(), "sampleChannel",
					"参加者一覧\nsample01@example.com\nsample02@example.com\n",
					"sampleThreadTimeStamp").Return(errors.New("sample error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			msr := mock_repository.NewMockSlackRepository(ctrl)
			mer := mock_repository.NewMockErrorRepository(ctrl)
			if tt.prepare != nil {
				tt.prepare(msr)
			}
			s := &slackResponseService{
				slackRepository: msr,
				errorRepository: mer,
			}
			if err := s.ReplyEmailList(context.Background(), tt.args.event, tt.args.emails); (err != nil) != tt.wantErr {
				t.Errorf("ReplyEmailList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_slackErrorResponseService_ReplyError(t *testing.T) {
	type args struct {
		event *slackevents.AppMentionEvent
		err   error
	}
	tests := []struct {
		name    string
		args    args
		prepare func(mer *mock_repository.MockErrorRepository, msr *mock_repository.MockSlackRepository)
		wantErr bool
	}{
		{
			name: "OK: err is ErrThreadNotFound",
			args: args{
				event: &slackevents.AppMentionEvent{
					Channel:         "sampleChannel",
					ThreadTimeStamp: "sampleThreadTimeStamp",
					User:            "sampleUser",
				},
				err: errors.New("thread_not_found"), //errors.New("user_not_found")
			},
			prepare: func(mer *mock_repository.MockErrorRepository, msr *mock_repository.MockSlackRepository) {
				gomock.InOrder(
					mer.EXPECT().ErrThreadNotFound(errors.New("thread_not_found")).Return(true),
					msr.EXPECT().PostEphemeral(gomock.Any(), "sampleChannel",
						"スレッドで呼び出してね:neko_namida:",
						"sampleThreadTimeStamp", "sampleUser").Return(nil),
				)
			},
		},
		{
			name: "OK: err is ErrUserNotFound",
			args: args{
				event: &slackevents.AppMentionEvent{
					Channel:         "sampleChannel",
					ThreadTimeStamp: "sampleThreadTimeStamp",
					User:            "sampleUser",
				},
				err: errors.New("user_not_found"),
			},
			prepare: func(mer *mock_repository.MockErrorRepository, msr *mock_repository.MockSlackRepository) {
				gomock.InOrder(
					mer.EXPECT().ErrThreadNotFound(errors.New("user_not_found")).Return(false),
					mer.EXPECT().ErrUserNotFound(errors.New("user_not_found")).Return(true),
					msr.EXPECT().PostMessage(gomock.Any(), "sampleChannel",
						"参加者はいないようです:neko_namida:",
						"sampleThreadTimeStamp").Return(nil),
				)
			},
		},
		{
			name: "NG: err is ErrThreadNotFound (error in slackRepository.PostEphemeral)",
			args: args{
				event: &slackevents.AppMentionEvent{
					Channel:         "sampleChannel",
					ThreadTimeStamp: "sampleThreadTimeStamp",
					User:            "sampleUser",
				},
				err: errors.New("thread_not_found"),
			},
			prepare: func(mer *mock_repository.MockErrorRepository, msr *mock_repository.MockSlackRepository) {
				gomock.InOrder(
					mer.EXPECT().ErrThreadNotFound(errors.New("thread_not_found")).Return(true),
					msr.EXPECT().PostEphemeral(gomock.Any(), "sampleChannel",
						"スレッドで呼び出してね:neko_namida:",
						"sampleThreadTimeStamp", "sampleUser").Return(errors.New("sample_error")),
				)
			},
			wantErr: true,
		},
		{
			name: "NG: err is ErrUserNotFound (error in slackRepository.PostMessage)",
			args: args{
				event: &slackevents.AppMentionEvent{
					Channel:         "sampleChannel",
					ThreadTimeStamp: "sampleThreadTimeStamp",
					User:            "sampleUser",
				},
				err: errors.New("user_not_found"),
			},
			prepare: func(mer *mock_repository.MockErrorRepository, msr *mock_repository.MockSlackRepository) {
				gomock.InOrder(
					mer.EXPECT().ErrThreadNotFound(errors.New("user_not_found")).Return(false),
					mer.EXPECT().ErrUserNotFound(errors.New("user_not_found")).Return(true),
					msr.EXPECT().PostMessage(gomock.Any(), "sampleChannel",
						"参加者はいないようです:neko_namida:",
						"sampleThreadTimeStamp").Return(errors.New("sample_error")),
				)
			},
			wantErr: true,
		},
		{
			name: "NG: undefined error",
			args: args{
				event: &slackevents.AppMentionEvent{
					Channel:         "sampleChannel",
					ThreadTimeStamp: "sampleThreadTimeStamp",
					User:            "sampleUser",
				},
				err: errors.New("undefined error"),
			},
			prepare: func(mer *mock_repository.MockErrorRepository, msr *mock_repository.MockSlackRepository) {
				gomock.InOrder(
					mer.EXPECT().ErrThreadNotFound(errors.New("undefined error")).Return(false),
					mer.EXPECT().ErrUserNotFound(errors.New("undefined error")).Return(false),
				)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			msr := mock_repository.NewMockSlackRepository(ctrl)
			mer := mock_repository.NewMockErrorRepository(ctrl)
			if tt.prepare != nil {
				tt.prepare(mer, msr)
			}
			s := &slackResponseService{
				slackRepository: msr,
				errorRepository: mer,
			}
			if err := s.ReplyError(context.Background(), tt.args.event, tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("ReplyError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_slackErrorResponseService_ReplyHelp(t *testing.T) {
	type args struct {
		event *slackevents.AppMentionEvent
	}
	tests := []struct {
		name    string
		args    args
		prepare func(mer *mock_repository.MockErrorRepository, msr *mock_repository.MockSlackRepository)
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				event: &slackevents.AppMentionEvent{
					Channel:         "sampleChannel",
					ThreadTimeStamp: "sampleThreadTimeStamp",
					User:            "sampleUser",
				},
			},
			prepare: func(mer *mock_repository.MockErrorRepository, msr *mock_repository.MockSlackRepository) {
				msr.EXPECT().PostEphemeral(
					gomock.Any(), "sampleChannel",
					"[使い方]\n"+
						"1. スレッドで `@Auriga :sanka:` のようにAurigaを呼び出し、リアクションを指定してください。\n"+
						"2. スレッドの開始メッセージに指定のリアクションをしたユーザのメールアドレス一覧を返します。\n"+
						"3. 結果をGoogleCalenderに貼り付けると一括招待できます！",
					"sampleThreadTimeStamp", "sampleUser").Return(nil)
			},
		},
		{
			name: "NG: error in slackRepository.PostEphemeral",
			args: args{
				event: &slackevents.AppMentionEvent{
					Channel:         "sampleChannel",
					ThreadTimeStamp: "sampleThreadTimeStamp",
					User:            "sampleUser",
				},
			},
			prepare: func(mer *mock_repository.MockErrorRepository, msr *mock_repository.MockSlackRepository) {
				msr.EXPECT().PostEphemeral(
					gomock.Any(), "sampleChannel",
					"[使い方]\n"+
						"1. スレッドで `@Auriga :sanka:` のようにAurigaを呼び出し、リアクションを指定してください。\n"+
						"2. スレッドの開始メッセージに指定のリアクションをしたユーザのメールアドレス一覧を返します。\n"+
						"3. 結果をGoogleCalenderに貼り付けると一括招待できます！",
					"sampleThreadTimeStamp", "sampleUser").Return(errors.New("sample error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			msr := mock_repository.NewMockSlackRepository(ctrl)
			mer := mock_repository.NewMockErrorRepository(ctrl)
			if tt.prepare != nil {
				tt.prepare(mer, msr)
			}
			s := &slackResponseService{
				slackRepository: msr,
				errorRepository: mer,
			}
			if err := s.ReplyHelp(context.Background(), tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("ReplyHelp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
