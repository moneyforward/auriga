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

package repository

import (
	"github.com/moneyforward/auriga/app/internal/domain/repository"
	"github.com/moneyforward/auriga/app/pkg/slack"
)

type factory struct {
	client slack.Client
}

func NewFactory(client slack.Client) *factory {
	return &factory{
		client: client,
	}
}

func (f *factory) SlackRepository() repository.SlackRepository {
	return newSlackRepository(f.client)
}

func (f *factory) ErrorRepository() repository.ErrorRepository {
	return newErrorRepository()
}
