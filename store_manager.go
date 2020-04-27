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
	aheadTimeout int
	provider     AccountProvider
	sync.RWMutex
}

// NewStoreManager 创建存储管理器
func NewStoreManager(aheadTimeout int, provider AccountProvider) *StoreManager {
	return &StoreManager{
		items:        make(map[string]*StoreItem, storeItemSize),
		aheadTimeout: aheadTimeout,
		provider:     provider,
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

	return item.token, nil
}

// RemoveTicket 移除ticket
func (m *StoreManager) RemoveTicket(appID string) {
	if item, ok := m.items[appID]; ok {
		item.SetTicket("", time.Now().Add(-1*time.Second))
	}
}
