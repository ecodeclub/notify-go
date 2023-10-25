// Copyright 2021 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sms

import (
	"context"
	"net/url"

	"github.com/ecodeclub/notify-go/pkg/notifier"
	"github.com/kevinburke/twilio-go"
	"github.com/pkg/errors"
)

type Config struct {
	AccountSID      string `json:"account_sid"`
	AuthToken       string `json:"auth_token"`
	FromPhoneNumber string `json:"from_phone_number"`
}

type twilioClient interface {
	SendMessage(from, to, body string, mediaURLs []*url.URL) (*twilio.Message, error)
}

type ChannelSmsImpl struct {
	client          twilioClient
	fromPhoneNumber string
}

type Content struct {
	Data string
}

func NewSmsChannel(c Config) *ChannelSmsImpl {
	client := twilio.NewClient(c.AccountSID, c.AuthToken, nil)
	return &ChannelSmsImpl{
		client:          client.Messages,
		fromPhoneNumber: c.FromPhoneNumber,
	}
}

func (sc *ChannelSmsImpl) Execute(ctx context.Context, deli notifier.Delivery) error {
	msgContent := sc.initSMSContent(deli.Content)

	for _, recv := range deli.Receivers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, err := sc.client.SendMessage(sc.fromPhoneNumber, recv.Phone, msgContent.Data, []*url.URL{})
			if err != nil {
				return errors.Wrapf(err, "failed to send message to phone number '%s' using Twilio", recv.Phone)
			}
		}
	}

	return nil
}

func (sc *ChannelSmsImpl) Name() string {
	return "sms"
}

func (sc *ChannelSmsImpl) initSMSContent(nc notifier.Content) Content {
	c := Content{
		Data: string(nc.Data),
	}
	return c
}
