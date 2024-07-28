package pkg

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Hook 是一个钩子函数。注意
// ctx 是一个有超时机制的 官方context.Context
// 必须处理超时的问题
type Hook func(ctx context.Context) error

func BuildCloseServerHook(servers ...Server) Hook {
	return func(ctx context.Context) error {
		wg := sync.WaitGroup{}
		doneCh := make(chan struct{})
		wg.Add(len(servers))
		//多协程的关掉多个server也行的
		for _, server := range servers {
			go func(server2 Server) {
				err := server2.Shutdown(ctx)
				if err != nil {
					fmt.Printf("server shutdown error: %v \n", err)
				}
				time.Sleep(time.Second)
				wg.Done()
			}(server)
		}
		//	再开一个协程判断是否所有服务链接断开
		go func() {
			wg.Wait()
			//只需有东西阻塞，不用传入什么
			doneCh <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			fmt.Printf("closing servers timeout \n")
			return ErrorHookTimeout
		case <-doneCh:
			fmt.Printf("close all servers \n")
			return nil
		}
	}
}
