package internal

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func getAppBadges() string {
	homeDir, _ := os.UserHomeDir()
	scriptFolder := filepath.Join(homeDir, ".config/uebersicht/simple-bar-server")
	script := `
		source .env/bin/activate \
		&& python src/app_badges.py
	`

	cmd := exec.Command("bash", "-c", script)
	cmd.Dir = scriptFolder

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error", err.Error(), string(stdout))
		return "{}"
	}
	return strings.TrimSpace(string(stdout))
}

func ScheduleGetAppBadges(appBadgesRefreshSec int64) {
	tick := time.Tick(time.Duration(appBadgesRefreshSec) * time.Second)
	for range tick {
		output := getAppBadges()
		log.Println(output)
	}
}
