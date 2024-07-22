package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type HandlerBasedOnTree struct {
	root *node
}

type node struct {
	path     string
	children []*node //这里只需要node地址更加节省内存

	//如果是叶子节点匹配后就能调用此方法
	handler handleFunc
}

func NewHandlerBasedOnTree() Handler {
	return &HandlerBasedOnTree{
		root: &node{},
	}
}

// Route 他又一个根节点里面是/是没什么用的
func (h *HandlerBasedOnTree) Route(method string, pattern string, handleFunc handleFunc) {

	//正常应该改成对应的error返回，但是为了兼容map懒得改了
	err := h.validatePattern(pattern)
	if err != nil {
		fmt.Printf("路由失败:%v\n", err)
	}
	//分割字符
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")

	cur := h.root

	for index, path := range paths {
		//当前节点遍历所有子节点查看是否符合path，不符合就创建节点，符合就往下找
		matchChild, ok := cur.findMatchChild(path)
		if ok {
			cur = matchChild
		} else {
			//	创建节点
			cur.createSubTree(paths[index:], handleFunc)
			return
		}
	}
	// 离开了循环，说明我们加入的是短路径，
	// 比如说我们先加入了 /order/detail
	// 再加入/order，那么会走到这里
	cur.handler = handleFunc

}

// 要创建后面的子树都要
func (n *node) createSubTree(paths []string, handleFunc handleFunc) {
	cur := n
	for _, path := range paths {
		nn := newNode(path)
		cur.children = append(cur.children, nn)
		cur = nn
	}
	cur.handler = handleFunc
}

func newNode(path string) *node {
	return &node{
		path:     path,
		children: make([]*node, 0, 2),
	}
}

func (n *node) findMatchChild(path string) (*node, bool) {
	for _, child := range n.children {
		//优先详细的
		if child.path == path &&
			child.path != "*" {
			return child, true
		}
		if child.path == path {
			return child, true
		}
	}
	return nil, false
}

func (h *HandlerBasedOnTree) ServerHTTP(c *Context) {
	handler, found := h.root.findRouter(c.R.URL.Path)
	if !found {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("not found"))
		return
	}
	handler(c)
}

func (n *node) findRouter(path string) (handleFunc, bool) {
	paths := strings.Split(strings.Trim(path, "/"), "/")
	cur := n
	for _, path := range paths {
		child, found := cur.findMatchChild(path)
		if !found {
			return nil, false
		}
		cur = child
	}
	if cur.handler == nil {
		return nil, false
	}
	return cur.handler, true
}

var ErrorInvalidRouterPattern = errors.New("invalid router pattern")

func (h *HandlerBasedOnTree) validatePattern(pattern string) error {
	// 校验 *，如果存在，必须在最后一个，并且它前面必须是/
	// 即我们只接受 /* 的存在，abc*这种是非法

	pos := strings.Index(pattern, "*")
	// 找到了 *
	if pos > 0 {
		// 必须是最后一个
		if pos != len(pattern)-1 {
			return ErrorInvalidRouterPattern
		}
		if pattern[pos-1] != '/' {
			return ErrorInvalidRouterPattern
		}
	}
	return nil
}
