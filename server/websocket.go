package server

import (
	_ "net/http/pprof"
)

var (
	lastLeaderboardUpdate = int64(0)
)
