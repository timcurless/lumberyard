package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"

	"github.com/timcurless/lumberyard/service"
)

func main() {
	var (
		httpaddr = flag.String("http.addr", ":8081", "HTTP Listen Address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s service.Service
	{
		s = service.NewCassandraService("127.0.0.1", "notused", "notused", "lumberyard")
		//s = service.NewInmemService()
		s = service.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = service.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "http", "addr", *httpaddr)
		errs <- http.ListenAndServe(*httpaddr, h)
	}()

	logger.Log("exit", <-errs)
}
