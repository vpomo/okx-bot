package httpcli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"okx-bot/exchange/logger"
	"strings"
	"time"
)

var Cli IHttpClient

func init() {
	Cli = NewDefaultHttpClient()
}

type DefaultHttpClient struct {
	cli     *http.Client
	timeout time.Duration
}

func NewDefaultHttpClient() *DefaultHttpClient {
	cli := new(DefaultHttpClient)
	cli.init()
	return cli
}

func (cli *DefaultHttpClient) init() {
	cli.timeout = 5 * time.Second
	cli.cli = &http.Client{
		Timeout: cli.timeout,
		Transport: &http.Transport{
			IdleConnTimeout:       time.Minute,
			TLSHandshakeTimeout:   cli.timeout,
			ResponseHeaderTimeout: cli.timeout,
			MaxConnsPerHost:       4,
			MaxIdleConnsPerHost:   4,
			MaxIdleConns:          8,
		},
	}
}

func (cli *DefaultHttpClient) SetTimeout(sec int64) {
	timeout := time.Duration(sec) * time.Second
	cli.timeout = timeout
	cli.cli.Timeout = timeout
	trans := cli.cli.Transport.(*http.Transport)
	trans.ResponseHeaderTimeout = timeout
	trans.TLSHandshakeTimeout = timeout
	trans.ExpectContinueTimeout = timeout
}

func (cli *DefaultHttpClient) SetProxy(proxy string) error {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		return err
	}
	trans := cli.cli.Transport.(*http.Transport)
	trans.Proxy = func(request *http.Request) (*url.URL, error) {
		return proxyUrl, nil
	}
	return nil
}

func (cli *DefaultHttpClient) DoRequest(method, rqUrl string, reqBody string, headers map[string]string) (data []byte, err error) {
	reqTimeoutCtx, _ := context.WithTimeout(context.TODO(), cli.timeout)
	req, _ := http.NewRequestWithContext(reqTimeoutCtx, method, rqUrl, strings.NewReader(reqBody))

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	resp, err := cli.cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("[DefaultHttpClient] close response body error:", err.Error())
		}
	}(resp.Body)

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %w", err)
	}

	if resp.StatusCode != 200 {
		return bodyData, errors.New(resp.Status)
	}

	return bodyData, nil
}
