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
	userWidgetIndex := query.Get("userWidgetIndex")

	wsMap.Add(target, userWidgetIndex, ws)

	if target == "" {
		slog.Warn("ws connection closed due to missing ?target=")
		ws.Close()
		return
	}

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			slog.Warn("Message received", "msgType", msgType, "error", err)
			break
		}
		slog.Info("Message received", "msgType", msgType, "msg", string(msg))
	}

	wsMap.Remove(target, userWidgetIndex, ws)
	slog.Info("ws connection closed", "target", target, "userWidgetIndex", userWidgetIndex)
}

func sendToWSClient(target string, userWidgetIndex string, payload any) {
	wsConns := wsMap.GetAll(target, userWidgetIndex)
	for _, wsConn := range wsConns {
		err := wsConn.WriteJSON(payload)
		if err != nil {
			slog.Warn("SendToWSClient", "error", err)
		}
	}
}

func CreateWebsocketRouter() *mux.Router {
	wsRouter := mux.NewRouter()
	wsRouter.HandleFunc("/", handleWebsocket).Methods("GET")
	return wsRouter
}
