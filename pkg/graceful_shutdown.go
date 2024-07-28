package pkg

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

var ErrorHookTimeout = errors.New("the hook timeout")

type GracefulShutdown struct {
	//还需要处理的请求数
	reqCnt int64
	//标识服务器开始拒绝外面的信息，也就是即将关闭
	closing int32
	//用通道来阻塞关闭
	zeroReqCnt chan struct{}
}

func NewGracefulShutdown() *GracefulShutdown {
	return &GracefulShutdown{
		zeroReqCnt: make(chan struct{}),
	}
}

// 注意这是个责任链关闭启动而已
func (g *GracefulShutdown) ShutdownFilterBuilder(next Filter) Filter {
	return func(c *Context) {
		//	开始拒绝其他请求
		cl := atomic.LoadInt32(&g.closing)
		//大于0就是服务要关闭了
		if cl > 0 {
			c.W.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		//处理当前服务
		atomic.AddInt64(&g.reqCnt, 1)
		next(c)
		//	处理完毕
		readyClose := atomic.AddInt64(&g.reqCnt, -1)
		// 已经开始关闭了，而且请求数为0，
		if cl > 0 && readyClose == 0 {
			g.zeroReqCnt <- struct{}{}
		}
	}
}

// RejectNewRequestAndWaiting 将会拒绝新的请求，并且等待处理中的请求
func (g *GracefulShutdown) RejectNewRequestAndWaiting(ctx context.Context) error {
	atomic.AddInt32(&g.closing, 1)
	//特殊情况关闭前就处理完了所有请求
	if atomic.LoadInt64(&g.reqCnt) == 0 {
		return nil
	}
	done := ctx.Done()
	select {
	case <-done:
		fmt.Println("超时了，还没等到所有请求执行完毕")
		return ErrorHookTimeout
	case <-g.zeroReqCnt:
		fmt.Println("全部请求处理完了")
	}
	return nil
}

// 关闭等待处理请求
func WaitForShutdown(hooks ...Hook) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, ShutdownSignals...)
	select {
	case sig := <-signals:
		fmt.Printf("get signal %s, application will shutdown \n", sig)
		//	十分钟还没结束就强行突出
		time.AfterFunc(time.Minute*10, func() {
			fmt.Printf("Shutdown gracefully timeout, application will shutdown immediately. ")
			os.Exit(1)
		})
		for _, hook := range hooks {
			//30s的倒计时，并且返回继承的上下文和对应的关闭函数
			ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
			err := hook(ctx)
			if err != nil {
				fmt.Printf("failed to run hook, err: %v \n", err)
			}
			cancelFunc()
		}
		os.Exit(0)
	}
}
