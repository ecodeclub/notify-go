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
	"net/http"

	"github.com/ecodeclub/ekit/bean/option"
	"github.com/ecodeclub/ekit/slice"
	"github.com/ecodeclub/notify-go/pkg/log"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"moul.io/http2curl"
)

type File resty.File

type Client struct {
	Service Resource
	Debug   bool
}

func WithDebug(debug bool) option.Option[Client] {
	return func(t *Client) {
		t.Debug = debug
	}
}

type Request struct {
	Header         map[string]string
	Query          map[string]string
	Body           any
	FormData       map[string]string
	PathParams     map[string]string
	UploadFiles    []File
	UploadFilePath map[string]string
	AuthScheme     string
	AuthToken      string
	BasicUserName  string
	BasicPassword  string
}

func NewClient(service Resource, opts ...option.Option[Client]) Client {
	c := Client{Service: service}
	option.Apply[Client](&c, opts...)
	return c
}

func (c Client) Ral(ctx context.Context, name string, req Request, respSuc any, respFail any) error {
	var lr = NewLogRecord()
	logger := log.FromContext(ctx)
	defer lr.Flush(logger)

	urlInfo, ok := slice.Find[Interface](c.Service.Interface, func(src Interface) bool {
		return src.Name == name
	})
	if !ok {
		return errors.New("[ral] 获取接口配置失败")
	}

	lr.Host = urlInfo.Host
	lr.Port = urlInfo.Port
	lr.Protocol = c.Service.Protocol
	lr.Url = urlInfo.URL
	lr.Method = urlInfo.Method

	rc := resty.New().EnableTrace().SetDebug(c.Debug).SetRetryCount(c.Service.Retry)
	rc.SetPreRequestHook(
		func(client *resty.Client, request *http.Request) error {
			command, err := http2curl.GetCurlCommand(request)
			logger.Info("", "curl", command)
			lr.CurlCmd = command.String()
			return err
		})
	rc.SetBaseURL(fmt.Sprintf("%s://%s:%s", c.Service.Protocol, urlInfo.Host, urlInfo.Port))

	client := rc.R().SetContext(ctx).
		SetHeaders(req.Header).
		SetQueryParams(req.Query).
		SetBody(req.Body).         // json data
		SetFormData(req.FormData). // form data
		SetResult(respSuc).
		SetError(respFail).
		SetPathParams(req.PathParams).
		SetAuthScheme(req.AuthScheme).SetAuthToken(req.AuthToken).
		SetBasicAuth(req.BasicUserName, req.BasicPassword)

	//通过文件流上传
	if len(req.UploadFiles) != 0 {
		for _, f := range req.UploadFiles {
			client.SetFileReader(f.ParamName, f.Name, f.Reader)
		}
	}

	//通过文件path上传
	if req.UploadFilePath != nil {
		client.SetFiles(req.UploadFilePath)
	}

	rsp, err := client.Execute(urlInfo.Method, urlInfo.URL)
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
