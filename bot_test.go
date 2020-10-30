package zero

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
	On("message", func(event Event, state State) bool {
		return event.RawMessage == "复读"
	}).Got("echo", "请输入复读内容", func(matcher *Matcher, event Event, state State) Response {
		Send(event, matcher.State["echo"])
		return SuccessResponse
	})
	Run(Option{
		Host:          "127.0.0.1",
		Port:          "6700",
		AccessToken:   "",
		NickName:      []string{"xcw", "镜华", "小仓唯"},
		CommandPrefix: "",
	})
	select {}
}
