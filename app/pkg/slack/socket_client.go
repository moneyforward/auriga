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

package slack

import (
	"log"
	"os"

	"github.com/slack-go/slack/socketmode"
)

type SocketClient interface {
	Events() chan socketmode.Event
	Debugf(format string, v ...interface{})
	Run() error
	Ack(req socketmode.Request, payload ...interface{})
}

type socketClient struct {
	*socketmode.Client
}

func NewSocketClient(client Client, debugMode bool) *socketClient {

	c := socketmode.New(client.GetClient(),
		socketmode.OptionDebug(debugMode),
		socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)

	return &socketClient{c}
}

func (c *socketClient) Events() chan socketmode.Event {
	return c.Client.Events
}

func (c *socketClient) Debugf(format string, v ...interface{}) {
	c.Client.Debugf(format, v...)
}

func (c *socketClient) Run() error {
	return c.Client.Run()
}

func (c *socketClient) Ack(req socketmode.Request, payload ...interface{}) {
	c.Client.Ack(req, payload...)
}
