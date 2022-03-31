package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func NewContext() context.Context {
	shutdownSignals := []os.Signal{os.Interrupt, syscall.SIGTERM}
	ctx, _ := signal.NotifyContext(context.Background(), shutdownSignals...)
	return ctx
}
