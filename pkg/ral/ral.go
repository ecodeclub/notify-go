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

package ral

import (
	"context"
	"fmt"

	"github.com/ecodeclub/notify-go/pkg/log"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type Client struct {
	Service Resource
}

func NewClient(service Resource) Client {
	return Client{service}
}

type Request struct {
	Header     map[string]string
	Query      map[string]string
	Body       any
	PathParams map[string]string
	AuthScheme string
	AuthToken  string
}

func (c Client) Ral(ctx context.Context, name string, req Request, respSucc any, respFail any) error {
	var lr = NewLogRecord()
	logger := log.FromContext(ctx)
	defer lr.Flush(logger)

	intf, ok := c.getUrl(name)
	if !ok {
		return errors.New("[ral] 获取接口配置失败")
	}

	lr.Host = intf.Host
	lr.Port = intf.Port
	lr.Protocol = c.Service.Protocol
	lr.Url = intf.URL
	lr.Method = intf.Method

	rc := resty.New().EnableTrace().SetRetryCount(c.Service.Retry)
	rc.SetBaseURL(fmt.Sprintf("%s://%s:%s", c.Service.Protocol, intf.Host, intf.Port))

	client := rc.R().SetContext(ctx).
		SetHeaders(req.Header).
		SetQueryParams(req.Query).
		SetBody(req.Body).
		SetResult(respSucc).
		SetError(respFail).
		SetPathParams(req.PathParams).
		SetAuthScheme(req.AuthScheme).SetAuthToken(req.AuthToken)

	rsp, err := client.Execute(intf.Method, intf.URL)
	lr.RspCode = rsp.StatusCode()
	lr.Error = rsp.Error()

	trace := rsp.Request.TraceInfo()
	lr.AddTimeCostPoint("total", trace.TotalTime)
	lr.AddTimeCostPoint("conn", trace.ConnTime)
	lr.AddTimeCostPoint("dns", trace.DNSLookup)
	lr.AddTimeCostPoint("server", trace.ServerTime)
	lr.AddTimeCostPoint("resp", trace.ResponseTime)
	lr.AddTimeCostPoint("tcp_conn", trace.TCPConnTime)
	lr.AddTimeCostPoint("tls_handshake", trace.TLSHandshake)

	return err
}

func (c Client) getUrl(name string) (Interface, bool) {
	intf := Interface{}
	for _, it := range c.Service.Interface {
		if it.Name == name {
			return it, true
		}
	}
	return intf, false
}
