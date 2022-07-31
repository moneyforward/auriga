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

//go:generate mockgen -source=slack.go -destination mock/slack.go
package repository

import (
	"context"

	"github.com/moneyforward/auriga/app/internal/model"
)

type SlackRepository interface {
	PostMessage(ctx context.Context, channelID, message, ts string) error
	PostEphemeral(ctx context.Context, channelID, message, ts, userID string) error
	GetParentMessage(ctx context.Context, channelID, ts string) (*model.SlackMessage, error)
	ListUsersEmail(ctx context.Context, userID []string) ([]*model.SlackUserEmail, error)
}
