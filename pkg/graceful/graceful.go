package graceful

import (
	"os"
	"os/signal"
)

func InitGraceful(cb []func()) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		for _, f := range cb {
			f()
		}
		close(c)
	}()
}
