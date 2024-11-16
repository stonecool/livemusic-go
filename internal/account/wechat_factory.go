package account

type WeChatFactory struct{}

//func (f *WeChatFactory) CreateAccount(cfg *config.AccountConfig) ICrawlAccount {
//	return &WeChatAccount{
//		Account: Account{
//			Category: "wechat",
//			State:    internal.AccountState_Uninitialized,
//			msgChan:  make(chan *internal.AsyncMessage),
//			done:     make(chan struct{}),
//		},
//		config: cfg,
//	}
//}
