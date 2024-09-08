package internal

import (
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

const defaultUserWidgetIndex = "0"

type WsMap struct {
	mu         sync.RWMutex
	wsConnsMap map[string][]*websocket.Conn
}

func (wsMap *WsMap) getKey(target string, userWidgetIndex string) string {
	if userWidgetIndex == "" {
		userWidgetIndex = defaultUserWidgetIndex
	}
	return target + "-" + userWidgetIndex
}

func (wsMap *WsMap) Add(target string, userWidgetIndex string, wsConn *websocket.Conn) {
	wsMap.mu.Lock()
	defer wsMap.mu.Unlock()

	key := wsMap.getKey(target, userWidgetIndex)
	wsConns, ok := wsMap.wsConnsMap[key]
	if !ok {
		wsConns = []*websocket.Conn{}
	}
	wsMap.wsConnsMap[key] = append(wsConns, wsConn)
}

func (wsMap *WsMap) Remove(target string, userWidgetIndex string, wsConn *websocket.Conn) {
	wsMap.mu.Lock()
	defer wsMap.mu.Unlock()

	key := wsMap.getKey(target, userWidgetIndex)
	if wsConns, ok := wsMap.wsConnsMap[key]; ok {
		wsMap.wsConnsMap[key] = slices.DeleteFunc(wsConns, func(conn *websocket.Conn) bool {
			return wsConn == conn
		})
	}
}

func (wsMap *WsMap) GetAll(target string, userWidgetIndex string) []*websocket.Conn {
	wsMap.mu.RLock()
	defer wsMap.mu.RUnlock()

	key := wsMap.getKey(target, userWidgetIndex)
	wsConns, ok := wsMap.wsConnsMap[key]
	if ok {
		return wsConns
	}
	return []*websocket.Conn{}
}
