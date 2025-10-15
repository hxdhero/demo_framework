package app

import (
	"context"
	"lls_api/pkg/rerr"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type App struct {
	servers []Server
	ctx     Context
}

func NewApp(servers []Server) (*App, error) {
	// 时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, rerr.Wrap(err)
	}
	time.Local = loc

	return &App{servers: servers}, nil
}

func (a *App) Run(ctx context.Context, initErrs []error) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	for i := range a.servers {
		wg.Add(1)
		go func(i int) {
			err := a.servers[i].Start()
			if err != nil {
				a.ctx.Log().Fatalf("Server start err: %v", err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	if len(initErrs) == 0 {
		a.ctx.Log().Info("服务都已启动")
	} else {
		a.ctx.Log().Info("服务都已启动,但是有如下初始化警告: ")
		for _, err := range initErrs {
			log.Default().Printf("HARAKIRI WARN!!! %s \n", rerr.ToString(err, false))
		}
	}

	// 等待信号
	select {
	case <-signals:
		// Received termination signal
		a.ctx.Log().Info("Received termination signal")
	case <-ctx.Done():
		// Context canceled
		a.ctx.Log().Info("Context canceled")
	}

	// 优雅退出服务
	for i := range a.servers {
		err := a.servers[i].Stop(ctx)
		if err != nil {
			a.ctx.Log().ErrorErr(rerr.WrapS(err, "Server stop err"))
		}
	}

	return nil
}
