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

func NewClient(name string) Client {
	c := Client{}
	for _, src := range service.Resources {
		if src.Name == name {
			c.Service = src
			break
		}
	}
	return c
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
