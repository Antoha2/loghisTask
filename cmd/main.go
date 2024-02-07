package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/antoha2/loghis/service"
	"github.com/antoha2/loghis/transport"
)

func main() {
	Run()
}

func Run() {

	// addrLog := config.GetLogConfig()
	// addrLogStr := fmt.Sprintf("http://%s:%d", addrLog.Host, addrLog.Port) //:8180

	// auth := authService.NewAuthService(addrAuthStr)

	logService := service.NewLogService()

	Tran := transport.NewMQ(logService)

	go Tran.Init()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	Tran.Stop()

}
