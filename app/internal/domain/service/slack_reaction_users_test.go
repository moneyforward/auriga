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
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	mock_repository "github.com/moneyforward/auriga/app/internal/domain/repository/mock"

	"github.com/moneyforward/auriga/app/internal/model"
)

func Test_slackReactionUsersService_ListUsersEmailByReaction(t *testing.T) {
	type args struct {
		channelID    string
		ts           string
		reactionName string
	}
	sampleMessage := &model.SlackMessage{
		ChannelID: "sampleCID",
		Reactions: []*model.SlackReaction{
			{
				Name:    "join",
				UserIDs: []string{"user01", "user02"},
			},
			{
				Name:    "reactionSample",
				UserIDs: []string{"user02", "user03"},
			},
		},
	}
	tests := []struct {
		name    string
		args    args
		prepare func(msr *mock_repository.MockSlackRepository)
		want    []*model.SlackUserEmail
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				channelID: "sampleCID", ts: "sampleTs", reactionName: "reactionSample",
			},
			prepare: func(msr *mock_repository.MockSlackRepository) {
				gomock.InOrder(
					msr.EXPECT().GetParentMessage(gomock.Any(), "sampleCID", "sampleTs").Return(
						sampleMessage, nil),
					msr.EXPECT().ListUsersEmail(gomock.Any(), []string{"user02", "user03"}).Return(
						[]*model.SlackUserEmail{
							{ID: "user02", Email: "user02@example.com"},
							{ID: "user03", Email: "user03@example.com"},
						}, nil),
				)
			},
			want: []*model.SlackUserEmail{
				{ID: "user02", Email: "user02@example.com"},
				{ID: "user03", Email: "user03@example.com"},
			},
		},
		{
			name: "NG: error in GetParentMessage",
			args: args{
				channelID: "sampleCID", ts: "sampleTs", reactionName: "reactionSample",
			},
			prepare: func(msr *mock_repository.MockSlackRepository) {
				gomock.InOrder(
					msr.EXPECT().GetParentMessage(gomock.Any(), "sampleCID", "sampleTs").Return(
						nil, errors.New("sample_error")),
				)
			},
			wantErr: true,
		},
		{
			name: "NG: error in ListUserEmail",
			args: args{
				channelID: "sampleCID", ts: "sampleTs", reactionName: "reactionSample",
			},
			prepare: func(msr *mock_repository.MockSlackRepository) {
				gomock.InOrder(
					msr.EXPECT().GetParentMessage(gomock.Any(), "sampleCID", "sampleTs").Return(
						sampleMessage, nil),
					msr.EXPECT().ListUsersEmail(gomock.Any(), []string{"user02", "user03"}).Return(
						nil, errors.New("sample_error")),
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
			if tt.prepare != nil {
				tt.prepare(msr)
			}
			s := &slackReactionUsersService{
				slackRepository: msr,
			}
			got, err := s.ListUsersEmailByReaction(context.Background(), tt.args.channelID, tt.args.ts, tt.args.reactionName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListUsersEmailByReaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListUsersEmailByReaction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slackReactionUsersService_getReactionUserIDs(t *testing.T) {
	type args struct {
		reactions    []*model.SlackReaction
		reactionName string
	}
	slackReactions := []*model.SlackReaction{
		{
			Name:    "join",
			UserIDs: []string{"user01", "user02"},
		},
		{
			Name:    "reactionSample",
			UserIDs: []string{"user02", "user03"},
		},
	}
	var noUsers []string
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "OK",
			args: args{
				reactions:    slackReactions,
				reactionName: "reactionSample",
			},
			want: []string{"user02", "user03"},
		},
		{
			name: "OK: no users to filter by reactionName",
			args: args{
				reactions:    slackReactions,
				reactionName: "sample",
			},
			want: noUsers,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			msr := mock_repository.NewMockSlackRepository(ctrl)
			s := &slackReactionUsersService{
				slackRepository: msr,
			}
			if got := s.getReactionUserIDs(context.Background(), tt.args.reactions, tt.args.reactionName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getReactionUserIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}
