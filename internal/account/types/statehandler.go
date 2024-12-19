package types

import "github.com/stonecool/livemusic-go/internal/message"

type StateHandler interface {
	GetState() message.AccountState
	SetState(message.AccountState)
	Transit(account Account, cmd message.AccountCmd)
	HandleError(account Account, err error)
}
