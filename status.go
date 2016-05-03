package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)



func statusService(addr string) {
	http.HandleFunc("/status", getStatus)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	status := make(map[string]interface{})

	sessionMut.Lock()
	// This can potentially be double the number of pending sessions, as each session has two keys, one for each side.
	status["startTime"] = rc.startTime
	status["uptimeSeconds"] = time.Since(rc.startTime) / time.Second
	status["numPendingSessionKeys"] = len(pendingSessions)
	status["numActiveSessions"] = len(activeSessions)
	sessionMut.Unlock()
	status["numConnections"] = atomic.LoadInt64(&numConnections)
	status["numProxies"] = atomic.LoadInt64(&numProxies)
	status["bytesProxied"] = atomic.LoadInt64(&bytesProxied)
	status["goVersion"] = runtime.Version()
	status["goOS"] = runtime.GOOS
	status["goArch"] = runtime.GOARCH
	status["goMaxProcs"] = runtime.GOMAXPROCS(-1)
	status["goNumRoutine"] = runtime.NumGoroutine()
	status["kbps10s1m5m15m30m60m"] = []int64{
		rc.rate(10/10) * 8 / 1000,
		rc.rate(60/10) * 8 / 1000,
		rc.rate(5*60/10) * 8 / 1000,
		rc.rate(15*60/10) * 8 / 1000,
		rc.rate(30*60/10) * 8 / 1000,
		rc.rate(60*60/10) * 8 / 1000,
	}
	status["options"] = map[string]interface{}{
		"network-timeout":  networkTimeout / time.Second,
		"ping-interval":    pingInterval / time.Second,
		"message-timeout":  messageTimeout / time.Second,
		"per-session-rate": sessionLimitBps,
		"global-rate":      globalLimitBps,
		"pools":            pools,
		"provided-by":      providedBy,
	}

	bs, err := json.MarshalIndent(status, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}
