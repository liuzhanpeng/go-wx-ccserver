package wxccserver

import (
	"testing"
)

type MockFetcher struct {
}

func NewMockFetcher() *MockFetcher {
	return &MockFetcher{}
}

func (f *MockFetcher) FetchToken(appID, appSecret string) (*TokenRes, error) {
	res := &TokenRes{
		AccessToken: "test-token1",
		ExpiresIn:   7200,
	}
	res.ErrCode = 0
	res.ErrMsg = ""

	return res, nil
}

func (f *MockFetcher) FetchTicket(accessToken string) (*TicketRes, error) {
	res := &TicketRes{
		Ticket:    "test-ticket1",
		ExpiresIn: 7200,
	}
	res.ErrCode = 0
	res.ErrMsg = ""

	return res, nil
}

func Test_StoreManager(t *testing.T) {
	manager := NewStoreManager(NewTomlAccountProvider("./cmd/accounts.toml"), NewMockFetcher(), 60)

	if len(manager.items) != 0 {
		t.Error("err1")
	}

	_, err := manager.Token("not exists")
	if err == nil {
		t.Error("err2")
	}

	manager.LoadAccounts()

	if len(manager.items) != 2 {
		t.Error("err3")
	}

	token, err := manager.Token("appidA")
	if err != nil {
		t.Error("err4")
	}

	manager.Token("appidA")

	if token != "test-token1" {
		t.Error("err5")
	}

	ticket, err := manager.Ticket("appidA")
	if err != nil {
		t.Error("err6")
	}

	if ticket != "test-ticket1" {
		t.Error("err7")
	}

}
