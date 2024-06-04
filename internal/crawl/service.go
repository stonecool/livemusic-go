package crawl

import (
	"log"
)

func newAccount(id int) (IAccount, error) {
	//account, err := GetAccountByID(id)
	//
	//if err != nil {
	//	return nil, err
	//}

	//if account.TemplateId == "WX" {
	//	return &wx.Account{Account: *account}, nil
	//}

	return nil, nil
}

func Test(id int) error {
	account, err := newAccount(id)
	if err != nil {
		return err
	}

	go start(account)
	return nil
}

func start(account IAccount) {
	log.Printf("Start account:%d\n", account.GetId())

	for {
		select {
		case cmd := <-account.GetChan():
			switch cmd.cmd {
			case CmdReady:
				if account.GetState() != StateInitial {
					log.Printf("state not initial")
					continue
				}

				ret, err := account.Login()
				if err != nil {
					log.Printf("error:%s", err)
					continue
				}

				if ret {
					account.SetState(StateReady)
				}
			case CmdRun:
				if account.GetState() != StateReady {
					log.Printf("state not ready")
					continue
				}

				account.SetState(StateRunning)
			case CmdSuspend:
				if account.GetState() != StateRunning {
					log.Printf("state not running")
					continue
				}

				account.SetState(StateReady)
			case CmdCrawl:
				if account.GetState() != StateRunning {
					log.Printf("state not running")
					continue
				}

				err := account.Crawl(cmd.instance)
				if err != nil {
					log.Printf("error:%s", err)
					// TODO
				}
			}
		}
	}
}

func startCrawl(instance *Instance) error {
	//if instance == nil {
	//	return fmt.Errorf("instance is nil")
	//}
	//
	//account, err := GetAccountByID(instance.AccountId)
	//if err != nil {
	//	log.Printf("%s", err)
	//	return err
	//}
	//
	//ticker := time.NewTicker(30 * time.Second)
	//defer ticker.Stop()

	//for {
	//	select {
	//	case <-ticker.C:
	//		if account.State == 1 {
	//			return nil
	//		}
	//
	//		account.GetChan() <- CmdRequest{
	//			cmd:      CmdCrawl,
	//			instance: instance,
	//		}
	//	}
	//}

	return nil
}
