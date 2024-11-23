package account

import (
	"github.com/stonecool/livemusic-go/internal/message"
)

type state int

const (
	stateNew state = iota
	stateInitialized
	stateNotLoggedIn
	stateReady
	stateRunning
	stateTerminated
)

func (s state) String() string {
	switch s {
	case stateNew:
		return "new"
	case stateInitialized:
		return "initialized"
	case stateNotLoggedIn:
		return "notLoggedIn"
	case stateReady:
		return "ready"
	case stateRunning:
		return "running"
	case stateTerminated:
		return "terminated"
	default:
		return "unknown"
	}
}

type stateManager interface {
	getNextState(currentState state, cmd message.CrawlCmd) state
	getErrorState(currentState state) state
	isValidTransition(from, to state) bool
}

type BaseStateManager struct{}

func (b *BaseStateManager) getNextState(state state, cmd message.CrawlCmd) state {
	switch state {
	case stateNew:
		if cmd == message.CrawlCmd_Initial {
			return stateInitialized
		}
	case stateReady:
		if cmd == message.CrawlCmd_Crawl {
			return stateRunning
		}
		//case stateRunning:
		//	if cmd == message.CrawlCmd_Stop {
		//		return stateReady
		//	}
	}
	return state
}

func (b *BaseStateManager) getErrorState(state state) state {
	switch state {
	case stateNew:
		return stateNew // 新建状态出错保持原状态
	case stateInitialized:
		return stateNew // 初始化失败回到新建状态
	case stateRunning:
		return stateReady // 运行出错回到就绪状态
	case stateTerminated:
		return stateTerminated // 终止状态不变
	default:
		return state
	}
}

func (b *BaseStateManager) isValidTransition(from, to state) bool {
	switch from {
	case stateNew:
		return to == stateInitialized
	case stateInitialized:
		return to == stateReady || to == stateNew
	case stateReady:
		return to == stateRunning || to == stateTerminated
	case stateRunning:
		return to == stateReady || to == stateTerminated
	case stateTerminated:
		return false
	default:
		return false
	}
}

// DefaultStateManager 需要登录的账号状态管理器
type DefaultStateManager struct {
	BaseStateManager
}

func (mgr *DefaultStateManager) getNextState(currentState state, cmd message.CrawlCmd) state {
	switch currentState {
	case stateInitialized:
		return stateNotLoggedIn
	case stateNotLoggedIn:
		if cmd == message.CrawlCmd_Login {
			return stateReady
		}
	}
	// 其他情况使用基类的默认实现
	return mgr.BaseStateManager.getNextState(currentState, cmd)
}

type NoLoginStateManager struct {
	BaseStateManager
}

func (mgr *NoLoginStateManager) getNextState(currentState state, cmd message.CrawlCmd) state {
	switch currentState {
	case stateInitialized:
		return stateReady // 直接进入就绪状态
	}
	// 其他情况使用基类的默认实现
	return mgr.BaseStateManager.getNextState(currentState, cmd)
}

func selectStateManager(category string) stateManager {
	switch category {
	case "wechat":
		return &DefaultStateManager{}
	case "noLogin":
		return &NoLoginStateManager{}
	default:
		return &DefaultStateManager{}
	}
}
