package web

import (
	"fmt"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

type Filter func(c *Context)
type FilterBuilder func(next Filter) Filter

type Server interface {
	Route(method string, pattern string, handleFunc func(ctx *Context))
	Start(address string) error
}

type sdkHttpServer struct {
	Name    string
	handler Handler
	root    Filter
}

type Handler interface {
	Route(method string, pattern string, handleFunc func(ctx *Context))
	ServeHTTP(c *Context)
}

type HandlerBaseOnMap struct {
	// 路由表
	handlers map[string]func(ctx *Context)
}

func NewHandlerBaseOnMap() *HandlerBaseOnMap {
	return &HandlerBaseOnMap{
		handlers: make(map[string]func(ctx *Context)),
	}
}

func (h *HandlerBaseOnMap) Route(method string, pattern string, handleFunc func(ctx *Context)) {
	key := fmt.Sprintf("%s#%s", method, pattern)
	h.handlers[key] = handleFunc
}

func (h *HandlerBaseOnMap) ServeHTTP(c *Context) {
	key := fmt.Sprintf("%s#%s", c.R.Method, c.R.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		handler(c)
	} else {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("404 - Not Found"))
	}
}

func logFilterBuilder(next Filter) Filter {
	return func(c *Context) {
		fmt.Println("Log Filter: Request received")
		next(c)
	}
}

func authFilterBuilder(next Filter) Filter {
	return func(c *Context) {
		fmt.Println("Auth Filter: Checking authentication")
		next(c)
	}
}

func NewServer(name string, builders ...FilterBuilder) Server {
	handler := NewHandlerBaseOnMap()
	var root Filter = func(c *Context) {
		handler.ServeHTTP(c)
	}

	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}
	return &sdkHttpServer{
		Name:    name,
		handler: handler,
		root:    root,
	}
}

func (s *sdkHttpServer) Route(method string, pattern string, handleFunc func(ctx *Context)) {
	s.handler.Route(method, pattern, handleFunc)
}

func (s *sdkHttpServer) Start(address string) error {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		c := NewContext(writer, request)
		s.root(c)
	})
	return http.ListenAndServe(address, nil)
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W: w,
		R: r,
	}
}

func main() {
	//server := NewServer("MyServer", logFilterBuilder, authFilterBuilder)
	//
	//server.Route("GET", "/hello", func(ctx *Context) {
	//	fmt.Fprintf(ctx.W, "Hello, world!")
	//})
	//
	//fmt.Println("Starting server at :8080")
	//if err := server.Start(":8080"); err != nil {
	//	fmt.Printf("Error starting server: %v\n", err)
	//}
}
