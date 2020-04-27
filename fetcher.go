package wxccserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// APIRes 接口通用响应
type APIRes struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// TokenRes 获取access token的接口的响应
type TokenRes struct {
	APIRes
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// TicketRes 获取jsapi_ticket的接口的响应
type TicketRes struct {
	APIRes
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

// Fetcher 是access token/jsapi ticket的请求器接口
type Fetcher interface {
	FetchToken(appID, appSecret string) (*TokenRes, error)
	FetchTicket(accessToken string) (*TicketRes, error)
}

// WxFetcher 微信公众号请求器;
// 调用公众号API获取相应的token和ticket
type WxFetcher struct {
	client *http.Client
}

// NewWxFetcher 创建请求器
func NewWxFetcher() *WxFetcher {
	return &WxFetcher{
		client: &http.Client{},
	}
}

// FetchToken 获取access token
func (f *WxFetcher) FetchToken(appID, appSecret string) (*TokenRes, error) {
	apiURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appID, appSecret)
	res, err := f.client.Get(apiURL)
	if err != nil {
		return nil, errors.Wrap(err, "请求access token发生网络错误")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求access token发生服务器错误, 服务返回状态码[%v]", res.StatusCode)
	}

	var tokenRes TokenRes
	err = json.NewDecoder(res.Body).Decode(&tokenRes)
	if err != nil {
		return nil, errors.Wrap(err, "access token解码失败")
	}

	if tokenRes.ErrCode != 0 {
		return nil, fmt.Errorf("请求access token发生业务异常:[%v]", tokenRes.ErrMsg)
	}

	return &tokenRes, nil
}

// FetchTicket 获取jsapi ticket
func (f *WxFetcher) FetchTicket(accessToken string) (*TicketRes, error) {
	apiURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", accessToken)
	res, err := f.client.Get(apiURL)
	if err != nil {
		return nil, errors.Wrap(err, "请求jsapi ticket发生网络错误")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求jsapi ticket发生服务器错误, 服务返回状态码[%v]", res.StatusCode)
	}

	var ticketRes TicketRes
	err = json.NewDecoder(res.Body).Decode(&ticketRes)
	if err != nil {
		return nil, errors.Wrap(err, "jsapi ticket解码失败")
	}

	if ticketRes.ErrCode != 0 {
		return nil, fmt.Errorf("请求jsapi ticket发生业务异常:[%v]", ticketRes.ErrMsg)
	}

	return &ticketRes, nil
}
