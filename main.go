package main

import (
	"simple-bar-server/internal"
)

const AppBadgesRefreshSec = 2

func main() {
	go internal.ScheduleGetAppBadges(AppBadgesRefreshSec)
	select {}
}
