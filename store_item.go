package wxccserver

import (
	"sync"
	"time"
)

// StoreItem 保存access token和jsapi ticket
type StoreItem struct {
	appID         string    // appid
	appSecret     string    // appsecret
	signatureKey  string    // 签名key
	token         string    // access token
	tokenExpired  time.Time // access token超时时间
	ticket        string    // jsapi ticket
	ticketExpired time.Time // jsapi ticket超时时间
	sync.RWMutex
}

// NewStoreItem 创建存储项
func NewStoreItem(appID, appSecret, signatureKey string) *StoreItem {
	return &StoreItem{
		appID:        appID,
		appSecret:    appSecret,
		signatureKey: signatureKey,
	}
}

// IsValidToken 判断access token是否有效; token不存在或超时均为无效
func (item *StoreItem) IsValidToken() bool {
	if item.token != "" && item.tokenExpired.After(time.Now()) {
		return true
	}

	return false
}

// Token 返回access token
func (item *StoreItem) Token() string {
	return item.token
}

// SetToken 设置access token
func (item *StoreItem) SetToken(token string, expired time.Time) {
	item.Lock()
	defer item.Unlock()

	item.token = token
	item.tokenExpired = expired
}

// IsValidTicket 判断jsapi ticket是否有效; ticket不存在或超时均为无效
func (item *StoreItem) IsValidTicket() bool {
	if item.ticket != "" && item.ticketExpired.After(time.Now()) {
		return true
	}

	return false
}

// Ticket 返回jsapi ticket
func (item *StoreItem) Ticket() string {
	return item.ticket
}

// SetTicket 设置jsapi ticket
func (item *StoreItem) SetTicket(ticket string, expired time.Time) {
	item.Lock()
	defer item.Unlock()

	item.ticket = ticket
	item.ticketExpired = expired
}
