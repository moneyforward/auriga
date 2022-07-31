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
	"reflect"
	"testing"

	"github.com/moneyforward/auriga/app/internal/model"
)

func TestNewSlackMentionedService(t *testing.T) {
	tests := []struct {
		name string
		want *slackMentionedService
	}{
		{
			name: "OK",
			want: &slackMentionedService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSlackMentionedService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSlackMentionedService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slackMentionedService_Parse(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want *model.MentionParseResult
	}{
		{
			name: "OK: when a reaction is specified",
			args: args{message: "@auriga :join:"},
			want: &model.MentionParseResult{
				Message:  "@auriga :join:",
				Reaction: "join",
			},
		},
		{
			name: "OK: help command",
			args: args{message: "@auriga help"},
			want: &model.MentionParseResult{
				Message: "@auriga help",
				Command: CommandHelp,
			},
		},
		{
			name: "NG: reaction is formatted incorrectly.",
			args: args{message: "@auriga :tmp"},
			want: &model.MentionParseResult{
				Message: "@auriga :tmp",
				Command: CommandHelp,
			},
		},
		{
			name: "NG: reaction is formatted incorrectly.",
			args: args{message: "@auriga tmp:"},
			want: &model.MentionParseResult{
				Message: "@auriga tmp:",
				Command: CommandHelp,
			},
		},
		{
			name: "NG: command that do not exist",
			args: args{message: "@auriga command"},
			want: &model.MentionParseResult{
				Message: "@auriga command",
				Command: CommandHelp,
			},
		},
		{
			name: "NG: no command is specified",
			args: args{message: "@auriga"},
			want: &model.MentionParseResult{
				Message: "@auriga",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &slackMentionedService{}
			if got := s.Parse(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
