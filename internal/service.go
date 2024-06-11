package internal

//func Test(id int) error {
//	crawl := GetCrawl(nil)
//
//	go start(crawl)
//	return nil
//}

//func start(crawl ICrawl) {
//	log.Printf("Start crawl:%d\n", crawl.GetId())
//
//	for {
//		select {
//		case cmd := <-crawl.GetChan():
//			switch cmd.cmd {
//			case CmdReady:
//				if crawl.GetState() != StateInitial {
//					log.Printf("state not initial")
//					continue
//				}
//
//				ret, err := crawl.Login()
//				if err != nil {
//					log.Printf("error:%s", err)
//					continue
//				}
//
//				if ret {
//					crawl.SetState(StateReady)
//				}
//			case CmdRun:
//				if crawl.GetState() != StateReady {
//					log.Printf("state not ready")
//					continue
//				}
//
//				crawl.SetState(StateRunning)
//			case CmdSuspend:
//				if crawl.GetState() != StateRunning {
//					log.Printf("state not running")
//					continue
//				}
//
//				crawl.SetState(StateReady)
//			case CmdCrawl:
//				if crawl.GetState() != StateRunning {
//					log.Printf("state not running")
//					continue
//				}
//
//				err := crawl.Crawl(cmd.instance)
//				if err != nil {
//					log.Printf("error:%s", err)
//					// TODO
//				}
//			}
//		}
//	}
//}

//func startCrawl(instance *Instance) error {
//if instance == nil {
//	return fmt.Errorf("instance is nil")
//}
//
//account, err := GetCrawlByID(instance.AccountId)
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

//	return nil
//}
