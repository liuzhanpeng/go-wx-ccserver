package wxccserver

import (
	"fmt"
	"sync"
	"time"
)

const storeItemSize = 32 // 初始存储项大小

// StoreManager 存储管理器; 用于管理各个公众号的access token和jsapi ticket
type StoreManager struct {
	items        map[string]*StoreItem
	provider     AccountProvider
	fetcher      Fetcher
	aheadTimeout int
	sync.RWMutex
}

// NewStoreManager 创建存储管理器
func NewStoreManager(provider AccountProvider, fetcher Fetcher, aheadTimeout int) *StoreManager {
	return &StoreManager{
		items:        make(map[string]*StoreItem, storeItemSize),
		provider:     provider,
		fetcher:      fetcher,
		aheadTimeout: aheadTimeout,
	}
}

// LoadAccounts 从配置加载公众号账号
func (m *StoreManager) LoadAccounts() error {
	accounts, err := m.provider.Obtain()
	if err != nil {
		return err
	}

	m.Lock()
	defer m.Unlock()

	m.items = make(map[string]*StoreItem, len(accounts))
	for _, account := range accounts {
		m.items[account.AppID] = NewStoreItem(account.AppID, account.AppSecret)
	}

	return nil
}

// Token 返回指定appID的access token
func (m *StoreManager) Token(appID string) (string, error) {
	item, ok := m.items[appID]
	if !ok {
		return "", fmt.Errorf("找不到appid为[%s]的公众号账号信息", appID)
	}

	if item.IsValidToken() {
		return item.token, nil
	}

	res, err := m.fetcher.FetchToken(item.appID, item.appSecret)
	if err != nil {
		return "", err
	}

	item.SetToken(res.AccessToken, time.Now().Add(time.Duration(res.ExpiresIn-m.aheadTimeout)*time.Second))

	return item.token, nil
}

// RemoveToken 移除token
func (m *StoreManager) RemoveToken(appID string) {
	if item, ok := m.items[appID]; ok {
		item.SetToken("", time.Now().Add(-1*time.Second))
	}
}

// Ticket 返回指定appID的jsapi ticket
func (m *StoreManager) Ticket(appID string) (string, error) {
	item, ok := m.items[appID]
	if !ok {
		return "", fmt.Errorf("找不到appid为[%s]的公众号账号信息", appID)
	}

	if item.IsValidTicket() {
		return item.ticket, nil
	}

	token, err := m.Token(appID)
	if err != nil {
		return "", err
	}

	res, err := m.fetcher.FetchTicket(token)
	if err != nil {
		return "", err
	}

	item.SetTicket(res.Ticket, time.Now().Add(time.Duration(res.ExpiresIn-m.aheadTimeout)*time.Second))

	return item.ticket, nil
}

// RemoveTicket 移除ticket
func (m *StoreManager) RemoveTicket(appID string) {
	if item, ok := m.items[appID]; ok {
		item.SetTicket("", time.Now().Add(-1*time.Second))
	}
}
