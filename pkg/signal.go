package pkg

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func OnInterrupt(cbs ...func()) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	// waiting for receiving signals
	<-sigs
	slog.Info("interrupt signal received")
	for _, cb := range cbs {
		cb()
	}
}
