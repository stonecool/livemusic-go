package account

import (
	"github.com/chromedp/chromedp"
	"github.com/stonecool/livemusic-go/internal/account/state"
	"github.com/stonecool/livemusic-go/internal/message"
)

type MockAccount struct {
	ID       int
	Category string
}

func (m *MockAccount) GetID() int {
	return m.ID
}

func (m *MockAccount) GetName() string {
	return "MockAccount"
}

func (m *MockAccount) getState() state.accountState {
	return state.stateNew
}

func (m *MockAccount) SetState(state state.accountState) {}

func (m *MockAccount) CheckLogin() chromedp.ActionFunc {
	return nil
}

func (m *MockAccount) WaitLogin() chromedp.ActionFunc {
	return nil
}

func (m *MockAccount) GetLoginURL() string {
	return ""
}

func (m *MockAccount) Login() error {
	return nil
}

func (m *MockAccount) GetQRCode([]byte) {}

func (m *MockAccount) GetQRCodeSelector() string {
	return ""
}

func (m *MockAccount) SaveCookies([]byte) error {
	return nil
}

func (m *MockAccount) GetCookies() []byte {
	return nil
}

func (m *MockAccount) GetLastURL() string {
	return ""
}

func (m *MockAccount) SetLastURL(url string) {}

func (m *MockAccount) IsAvailable() bool {
	return true
}

func (m *MockAccount) GetMsgChan() chan *message.AsyncMessage {
	return make(chan *message.AsyncMessage)
}

func (m *MockAccount) GetCategory() string {
	return m.Category
}
