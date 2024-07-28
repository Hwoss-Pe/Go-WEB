package v1

import (
	"net/http"
	"sync"
)

type Routable interface {
	Route(method string, pattern string, handleFunc handleFunc) error
}
type HandlerBaseOnMap struct {
	//key是对应的请求方式加路径
	handlers sync.Map
}

// 启动的时候进行路由
func (h *HandlerBaseOnMap) Route(
	method string, pattern string,
	handleFunc handleFunc) error {

	key := h.key(method, pattern)
	h.handlers.Store(key, handleFunc)
	return nil
}

// ServerHTTP http.Handler 下面方法是让map实现了这个Handler，因此那边用的不再是HandleFunc
func (h *HandlerBaseOnMap) ServerHTTP(ctx *Context) {
	key := h.key(ctx.R.Method, ctx.R.URL.Path)
	handler, ok := h.handlers.Load(key)
	if !ok {
		ctx.W.WriteHeader(http.StatusNotFound)
		_, _ = ctx.W.Write([]byte("not any router match"))
		return
	}
	//对于syncMap的东西本质上都是类似一个泛型，因此要断言一下类型
	handler.(handleFunc)(ctx)
}
func (h *HandlerBaseOnMap) key(method string, pattern string) string {
	return method + "#" + pattern
}

// 进行断言HandlerBaseOnMap一定实现了Handler接口
//这里已经被更改如果需要实现接口需要实现ServerHTTP并且注释树里面的
//type Handler interface {
//	ServeHTTP(c *Context)
//	Routable
//}
//
//
//var _ Handler = &HandlerBaseOnMap{}

func NewHandlerBaseOnMap() *HandlerBaseOnMap {
	return &HandlerBaseOnMap{
		//这里就是路由修改
	}
}
