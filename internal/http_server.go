package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
)

var (
	aerospaceKinds   = []string{"spaces", "windows"}
	aerospaceActions = []string{"refresh"}
)

func handleAerospace(res http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	kind := vars["kind"]
	if !slices.Contains(aerospaceKinds, kind) {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "Unknown aerospace kind %s", kind)
		return nil
	}

	action := vars["action"]
	if !slices.Contains(aerospaceActions, action) {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "Unknown aerospace action %s", action)
		return nil
	}

	data := map[string]any{}
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil && err != io.EOF {
		return err
	}

	payload := map[string]any{"action": action, "data": data}
	sendToWSClient(kind, "", payload)
	res.WriteHeader(http.StatusOK)
	return nil
}

var widgetKinds = []string{
	"app-badges",
	"battery",
	"browser-track",
	"crypto",
	"date-display",
	"keyboard",
	"mic",
	"mpd",
	"music",
	"netstats",
	"cpu",
	"sound",
	"spotify",
	"stock",
	"time",
	"user-widget",
	"viscosity-vpn",
	"weather",
	"wifi",
	"zoom",
}
var widgetActions = []string{"toggle", "enable", "disable", "refresh"}

func handleWidget(res http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	kind := vars["kind"]
	if !slices.Contains(widgetKinds, kind) {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "Unknown aerospace kind %s", kind)
		return nil
	}

	action := vars["action"]
	if !slices.Contains(widgetActions, action) {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "Unknown aerospace action %s", action)
		return nil
	}

	payload := map[string]string{"action": action}
	sendToWSClient(kind, vars["userWidgetIndex"], payload)
	res.WriteHeader(http.StatusOK)
	return nil
}

func CreateHTTPRouter() *mux.Router {
	httpRouter := mux.NewRouter()
	httpRouter.HandleFunc("/aerospace/{kind}/{action}", loggingMiddleware(handleAerospace)).Methods("POST")
	httpRouter.HandleFunc("/widget/{kind}/{action}/{userWidgetIndex}", loggingMiddleware(handleWidget)).Methods("POST")
	return httpRouter
}
