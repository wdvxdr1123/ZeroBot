package ZeroBot

import (
	"net/http"
	"net/http/pprof"
	"testing"
)

const (
	pprofAddr string = ":7890"
)

func StartHTTPDebuger() {
	pprofHandler := http.NewServeMux()
	pprofHandler.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	server := &http.Server{Addr: pprofAddr, Handler: pprofHandler}
	go server.ListenAndServe()
}

func TestRun(t *testing.T) {
	go StartHTTPDebuger()
	On(func(event Event) bool {
		if tp, ok := event["post_type"]; !ok || tp.String() != "message" {
			return false
		}
		return event["raw_message"].Str == "复读"
	}).Got("echo", "请输入复读内容", func(event Event, matcher *Matcher) Response {
		Send(event, matcher.State["echo"])
		return SuccessResponse
	})
	Run("ws://127.0.0.1:6700", "")
	select {}
}