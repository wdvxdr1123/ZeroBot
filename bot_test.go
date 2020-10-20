package ZeroBot

import "testing"

func TestRun(t *testing.T) {
	Run("ws://127.0.0.1:6700", "")
	select {}
}
