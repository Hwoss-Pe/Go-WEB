package v1

import "net/http"

type Routable interface {
	Route(method string, pattern string, handleFunc func(ctx *Context))
}
type HandlerBaseOnMap struct {
	//key是对应的请求方式加路径
	handlers map[string]func(ctx *Context)
}

type Handler interface {
	ServerHTTP(ctx *Context)
	Routable
}

// 启动的时候进行路由
func (h *HandlerBaseOnMap) Route(
	method string, pattern string,
	handleFunc func(ctx *Context)) {

	key := h.key(method, pattern)
	h.handlers[key] = handleFunc
}

// ServerHTTP http.Handler 下面方法是让map实现了这个Handler，因此那边用的不再是HandleFunc
func (h *HandlerBaseOnMap) ServerHTTP(ctx *Context) {
	key := h.key(ctx.R.Method, ctx.R.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		//		注册过
		handler(NewContext(ctx.W, ctx.R))
	} else {
		ctx.W.WriteHeader(http.StatusNotFound)
		ctx.W.Write([]byte("Not Found"))
	}
	h.handlers[key] = func(ctx *Context) {

	}
}
func (h *HandlerBaseOnMap) key(method string, pattern string) string {
	return method + "#" + pattern
}

// 进行断言HandlerBaseOnMap一定实现了Handler接口
var _ Handler = &HandlerBaseOnMap{}

func NewHandlerBaseOnMap() Handler {
	return &HandlerBaseOnMap{
		//这里就是路由修改
		handlers: make(map[string]func(c *Context), 128),
	}
}
