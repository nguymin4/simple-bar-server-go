package internal

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var wsMap = &WsMap{
	wsConnsMap: make(map[string][]*websocket.Conn),
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

func handleWebsocket(res http.ResponseWriter, req *http.Request) {
	ws, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		slog.Error("Failed to upgrade HTTP to Websocket", "error", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	query := req.URL.Query()
	target := query.Get("target")
	if target == "" {
		slog.Warn("ws connection closed due to missing ?target=")
		ws.Close()
		return
	}

	userWidgetIndex := query.Get("userWidgetIndex")
	wsMap.Add(target, userWidgetIndex, ws)
	slog.Info("ws connection opened", "target", target, "userWidgetIndex", userWidgetIndex)

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		slog.Info("Message received", "msgType", msgType, "msg", string(msg))
	}

	wsMap.Remove(target, userWidgetIndex, ws)
	slog.Info("ws connection closed", "target", target, "userWidgetIndex", userWidgetIndex, "error", err)
}

func sendToWSClient(target string, userWidgetIndex string, payload any) {
	wsConns := wsMap.GetAll(target, userWidgetIndex)
	if len(wsConns) == 0 {
		slog.Warn("no ws connection found", "target", target, "userWidgetIndex", userWidgetIndex, "payload", payload)
	}

	for _, wsConn := range wsConns {
		err := wsConn.WriteJSON(payload)
		if err != nil {
			slog.Warn("sendToWSClient", "error", err)
		}
	}
}

func CreateWebsocketRouter() *mux.Router {
	wsRouter := mux.NewRouter()
	wsRouter.HandleFunc("/", handleWebsocket).Methods("GET")
	return wsRouter
}
