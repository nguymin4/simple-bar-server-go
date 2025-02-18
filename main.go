package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
	"simple-bar-server/internal"
)

const (
	appBadgesRefreshSec = 2
	httpPort            = 7776
	wsPort              = 7777
)

func startServer(router *mux.Router, port int) {
	address := fmt.Sprintf(":%v", port)
	slog.Info("Started listening at " + address)

	err := http.ListenAndServe(address, router)
	if err != nil {
		slog.Error("Error listening for server", "error", err)
	}
}

func refreshUebersicht() {
	cmd := exec.Command("/usr/bin/osascript", "-e", "tell application id \"tracesOf.Uebersicht\" to refresh")
	err := cmd.Run()
	if err != nil {
		slog.Warn("Failed to refresh Uebersicht", "error", err)
	} else {
		slog.Info("Refreshed Uebersicht")
	}
}

func main() {
	go startServer(internal.CreateHTTPRouter(), httpPort)
	go startServer(internal.CreateWebsocketRouter(), wsPort)
	go internal.ScheduleGetAppBadges(appBadgesRefreshSec)
	go func() {
		time.Sleep(time.Duration(1) * time.Second)
		refreshUebersicht()
	}()
	select {}
}
