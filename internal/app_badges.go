package internal

import (
	"encoding/json"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func getAppBadges() (string, error) {
	homeDir, _ := os.UserHomeDir()
	scriptFolder := filepath.Join(homeDir, ".config/uebersicht/simple-bar-server-go/python")
	script := `
		source .env/bin/activate \
		&& python app_badges.py
	`

	cmd := exec.Command("bash", "-c", script)
	cmd.Dir = scriptFolder

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return string(stdout), err
	}
	return strings.TrimSpace(string(stdout)), nil
}

func ScheduleGetAppBadges(appBadgesRefreshSec int64) {
	tick := time.Tick(time.Duration(appBadgesRefreshSec) * time.Second)
	for range tick {
		output, err := getAppBadges()
		if err != nil {
			slog.Warn("Failed to get app badges", "error", err, "output", output)
			continue
		}

		data := map[string]any{}
		err = json.Unmarshal([]byte(output), &data)
		if err != nil {
			slog.Error("Failed to parse app badges", "error", err)
			continue
		}

		payload := map[string]any{"action": "refresh", "data": data}
		sendToWSClient("app-badges", "", payload)
		slog.Info("Finished updating app badges")
	}
}
