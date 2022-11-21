package signals

import (
	"context"
	"os"
	"os/signal"
)

var onlyOneSignalHandler = make(chan struct{})

// SetupSignalHandler 设置 stopCh 信号处理
func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	gracefulStopCh := make(chan os.Signal, 2)
	signal.Notify(gracefulStopCh, shutdownSignals...)
	go func() {
		// waiting for os signal to stop the program
		<-gracefulStopCh
		close(stop)
		<-gracefulStopCh
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

// GracefulStopWithContext 获取 ctx 信号处理
func GracefulStopWithContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	gracefulStopCh := make(chan os.Signal, 2)
	signal.Notify(gracefulStopCh, shutdownSignals...)
	go func() {
		// waiting for os signal to stop the program
		<-gracefulStopCh
		cancel()
		<-gracefulStopCh
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx
}
